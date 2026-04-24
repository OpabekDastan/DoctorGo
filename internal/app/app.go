package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"doctor_go/internal/config"
	"doctor_go/internal/handler"
	"doctor_go/internal/middleware"
	"doctor_go/internal/model"
	"doctor_go/internal/repository/postgres"
	"doctor_go/internal/service"
	"doctor_go/internal/worker"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := postgres.New(cfg.DSN())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := postgres.AutoMigrate(db, "migrations"); err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// --- Repositories ---
	authRepo        := postgres.NewAuthRepository(db)
	doctorRepo      := postgres.NewDoctorRepository(db)
	patientRepo     := postgres.NewPatientRepository(db)
	scheduleRepo    := postgres.NewScheduleRepository(db)
	appointmentRepo := postgres.NewAppointmentRepository(db)
	dashboardRepo   := postgres.NewDashboardRepository(db)
	notifRepo       := postgres.NewNotificationRepository(db)

	// --- Services ---
	authService        := service.NewAuthService(authRepo, cfg.JWTSecret, cfg.AccessTokenTTLMin)
	doctorService      := service.NewDoctorService(doctorRepo)
	patientService     := service.NewPatientService(patientRepo)
	scheduleService    := service.NewScheduleService(scheduleRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	dashboardService   := service.NewDashboardService(dashboardRepo)

	// --- Background Worker ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reminderWorker := worker.NewReminderWorker(appointmentService, notifRepo)
	reminderWorker.Start(ctx)

	// Graceful shutdown: отменяем ctx при SIGTERM / SIGINT
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down...")
		cancel()
	}()

	// --- Handlers ---
	healthHandler           := handler.NewHealthHandler()
	authHandler             := handler.NewAuthHandler(authService)
	adminDoctorHandler      := handler.NewAdminDoctorHandler(doctorService)
	adminPatientHandler     := handler.NewAdminPatientHandler(patientService)
	adminAppointmentHandler := handler.NewAdminAppointmentHandler(appointmentService)
	adminDashboardHandler   := handler.NewAdminDashboardHandler(dashboardService)
	doctorScheduleHandler   := handler.NewDoctorScheduleHandler(scheduleService)
	doctorAppointmentHandler := handler.NewDoctorAppointmentHandler(appointmentService)

	// --- Routes ---
	r.GET("/health", healthHandler.Health)

	api := r.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)

		authGroup := api.Group("")
		authGroup.Use(middleware.Auth(cfg.JWTSecret))
		{
			authGroup.GET("/me", authHandler.Me)

			adminGroup := authGroup.Group("/admin")
			adminGroup.Use(middleware.RequireRole(model.RoleAdmin))
			{
				// Dashboard
				adminGroup.GET("/dashboard", adminDashboardHandler.GetStats)

				// Doctors
				adminGroup.GET("/doctors", adminDoctorHandler.List)
				adminGroup.POST("/doctors", adminDoctorHandler.Create)
				adminGroup.GET("/doctors/:id", adminDoctorHandler.Get)
				adminGroup.PUT("/doctors/:id", adminDoctorHandler.Update)
				adminGroup.DELETE("/doctors/:id", adminDoctorHandler.Delete)

				// Patients
				adminGroup.GET("/patients", adminPatientHandler.List)
				adminGroup.POST("/patients", adminPatientHandler.Create)
				adminGroup.GET("/patients/:id", adminPatientHandler.Get)
				adminGroup.PUT("/patients/:id", adminPatientHandler.Update)
				adminGroup.DELETE("/patients/:id", adminPatientHandler.Delete)

				// Appointments
				adminGroup.GET("/appointments", adminAppointmentHandler.List)
				adminGroup.POST("/appointments", adminAppointmentHandler.Create)
				adminGroup.GET("/appointments/:id", adminAppointmentHandler.Get)
				adminGroup.PUT("/appointments/:id", adminAppointmentHandler.Update)
				adminGroup.DELETE("/appointments/:id", adminAppointmentHandler.Delete)
			}

			doctorGroup := authGroup.Group("/doctor")
			doctorGroup.Use(middleware.RequireRole(model.RoleDoctor))
			{
				doctorGroup.GET("/schedules", doctorScheduleHandler.List)
				doctorGroup.POST("/schedules", doctorScheduleHandler.Create)
				doctorGroup.GET("/schedules/:id", doctorScheduleHandler.Get)
				doctorGroup.PUT("/schedules/:id", doctorScheduleHandler.Update)
				doctorGroup.DELETE("/schedules/:id", doctorScheduleHandler.Delete)

				doctorGroup.GET("/appointments", doctorAppointmentHandler.List)
				doctorGroup.GET("/appointments/:id", doctorAppointmentHandler.Get)
				doctorGroup.PATCH("/appointments/:id/status", doctorAppointmentHandler.SetStatus)
			}
		}
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("server started on %s", addr)
	return r.Run(addr)
}
