package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"doctor_go/internal/model"
	"doctor_go/internal/repository/postgres"
	"doctor_go/internal/service"
)

type ReminderWorker struct {
	appointmentSvc *service.AppointmentService
	notifRepo      *postgres.NotificationRepository
}

func NewReminderWorker(
	appointmentSvc *service.AppointmentService,
	notifRepo *postgres.NotificationRepository,
) *ReminderWorker {
	return &ReminderWorker{
		appointmentSvc: appointmentSvc,
		notifRepo:      notifRepo,
	}
}

// Start запускает воркер в фоновой горутине.
// Останавливается когда ctx отменяется (например, при shutdown сервера).
func (w *ReminderWorker) Start(ctx context.Context) {
	// Для демо — каждую минуту. В проде можно поставить time.Hour.
	ticker := time.NewTicker(1 * time.Minute)
	log.Println("[Worker] Reminder worker started")

	go func() {
		for {
			select {
			case <-ticker.C:
				w.processReminders()
			case <-ctx.Done():
				ticker.Stop()
				log.Println("[Worker] Reminder worker stopped")
				return
			}
		}
	}()
}

func (w *ReminderWorker) processReminders() {
	log.Println("[Worker] Scanning for upcoming appointments...")

	appointments, err := w.appointmentSvc.ListAdmin()
	if err != nil {
		log.Printf("[Worker] ERROR fetching appointments: %v", err)
		return
	}

	now := time.Now()
	sent := 0

	for _, appt := range appointments {
		// Напоминаем за 24 часа до записи
		hoursUntil := appt.StartAt.Sub(now).Hours()
		if appt.Status == model.AppointmentScheduled && hoursUntil > 0 && hoursUntil <= 24 {
			msg := fmt.Sprintf(
				"Reminder: appointment #%d for %s %s scheduled at %s",
				appt.ID,
				appt.PatientFirstName,
				appt.PatientLastName,
				appt.StartAt.Format("2006-01-02 15:04"),
			)

			if _, err := w.notifRepo.Create(appt.ID, msg); err != nil {
				log.Printf("[Worker] ERROR saving notification for appointment %d: %v", appt.ID, err)
				continue
			}

			log.Printf("[Worker] 🔔 Reminder sent: %s", msg)
			sent++
		}
	}

	log.Printf("[Worker] Done. Reminders sent: %d", sent)
}
