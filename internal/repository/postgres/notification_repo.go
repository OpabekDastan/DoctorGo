package postgres

import (
	"doctor_go/internal/model"

	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(appointmentID int64, message string) (*model.Notification, error) {
	var notif model.Notification
	query := `
		INSERT INTO notifications (appointment_id, message)
		VALUES ($1, $2)
		RETURNING id, appointment_id, message, sent_at
	`
	if err := r.db.Get(&notif, query, appointmentID, message); err != nil {
		return nil, err
	}
	return &notif, nil
}

func (r *NotificationRepository) ListByAppointment(appointmentID int64) ([]model.Notification, error) {
	items := []model.Notification{}
	query := `SELECT id, appointment_id, message, sent_at FROM notifications WHERE appointment_id = $1 ORDER BY sent_at DESC`
	if err := r.db.Select(&items, query, appointmentID); err != nil {
		return nil, err
	}
	return items, nil
}
