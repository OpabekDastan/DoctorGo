package model

import "time"

type Notification struct {
	ID            int64     `db:"id"             json:"id"`
	AppointmentID int64     `db:"appointment_id" json:"appointment_id"`
	Message       string    `db:"message"        json:"message"`
	SentAt        time.Time `db:"sent_at"        json:"sent_at"`
}
