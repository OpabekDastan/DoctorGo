package model

import "time"

const (
	RoleAdmin  = "ADMIN"
	RoleDoctor = "DOCTOR"
)

const (
	AppointmentScheduled = "SCHEDULED"
	AppointmentConfirmed = "CONFIRMED"
	AppointmentCompleted = "COMPLETED"
	AppointmentCancelled = "CANCELLED"
	AppointmentNoShow    = "NO_SHOW"
)

type User struct {
	ID           int64     `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Doctor struct {
	ID             int64     `db:"id" json:"id"`
	Email          string    `db:"email" json:"email"`
	Role           string    `db:"role" json:"role"`
	FirstName      string    `db:"first_name" json:"first_name"`
	LastName       string    `db:"last_name" json:"last_name"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	Specialization string    `db:"specialization" json:"specialization"`
	Phone          string    `db:"phone" json:"phone"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type Patient struct {
	ID        int64     `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	BirthDate *string   `db:"birth_date" json:"birth_date,omitempty"`
	Gender    string    `db:"gender" json:"gender"`
	Phone     string    `db:"phone" json:"phone"`
	Email     string    `db:"email" json:"email"`
	Address   string    `db:"address" json:"address"`
	Comment   string    `db:"comment" json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DoctorSchedule struct {
	ID          int64  `db:"id" json:"id"`
	DoctorID    int64  `db:"doctor_id" json:"doctor_id"`
	Weekday     int    `db:"weekday" json:"weekday"`
	StartTime   string `db:"start_time" json:"start_time"`
	EndTime     string `db:"end_time" json:"end_time"`
	SlotMinutes int    `db:"slot_minutes" json:"slot_minutes"`
}

type Appointment struct {
	ID        int64     `db:"id" json:"id"`
	PatientID int64     `db:"patient_id" json:"patient_id"`
	DoctorID  int64     `db:"doctor_id" json:"doctor_id"`
	StartAt   time.Time `db:"start_at" json:"start_at"`
	EndAt     time.Time `db:"end_at" json:"end_at"`
	Status    string    `db:"status" json:"status"`
	Reason    string    `db:"reason" json:"reason"`
	Comment   string    `db:"comment" json:"comment"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type AppointmentWithNames struct {
	Appointment
	PatientFirstName string `db:"patient_first_name" json:"patient_first_name"`
	PatientLastName  string `db:"patient_last_name" json:"patient_last_name"`
	DoctorFirstName  string `db:"doctor_first_name" json:"doctor_first_name"`
	DoctorLastName   string `db:"doctor_last_name" json:"doctor_last_name"`
	DoctorEmail      string `db:"doctor_email" json:"doctor_email"`
}
