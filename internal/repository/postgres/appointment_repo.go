package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"doctor_go/internal/model"
)

type AppointmentRepository struct {
	db *sqlx.DB
}

func NewAppointmentRepository(db *sqlx.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) ListAdmin() ([]model.AppointmentWithNames, error) {
	items := []model.AppointmentWithNames{}
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.start_at, a.end_at, a.status, a.reason, a.comment, a.created_by, a.created_at,
		       p.first_name AS patient_first_name, p.last_name AS patient_last_name,
		       u.first_name AS doctor_first_name, u.last_name AS doctor_last_name, u.email AS doctor_email
		FROM appointments a
		JOIN patients p ON p.id = a.patient_id
		JOIN users u ON u.id = a.doctor_id
		ORDER BY a.start_at DESC
	`
	if err := r.db.Select(&items, query); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AppointmentRepository) ListByDoctor(doctorID int64) ([]model.AppointmentWithNames, error) {
	items := []model.AppointmentWithNames{}
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.start_at, a.end_at, a.status, a.reason, a.comment, a.created_by, a.created_at,
		       p.first_name AS patient_first_name, p.last_name AS patient_last_name,
		       u.first_name AS doctor_first_name, u.last_name AS doctor_last_name, u.email AS doctor_email
		FROM appointments a
		JOIN patients p ON p.id = a.patient_id
		JOIN users u ON u.id = a.doctor_id
		WHERE a.doctor_id = $1
		ORDER BY a.start_at DESC
	`
	if err := r.db.Select(&items, query, doctorID); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AppointmentRepository) GetAdminByID(id int64) (*model.AppointmentWithNames, error) {
	var item model.AppointmentWithNames
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.start_at, a.end_at, a.status, a.reason, a.comment, a.created_by, a.created_at,
		       p.first_name AS patient_first_name, p.last_name AS patient_last_name,
		       u.first_name AS doctor_first_name, u.last_name AS doctor_last_name, u.email AS doctor_email
		FROM appointments a
		JOIN patients p ON p.id = a.patient_id
		JOIN users u ON u.id = a.doctor_id
		WHERE a.id = $1
	`
	if err := r.db.Get(&item, query, id); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AppointmentRepository) GetDoctorByID(id, doctorID int64) (*model.AppointmentWithNames, error) {
	var item model.AppointmentWithNames
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.start_at, a.end_at, a.status, a.reason, a.comment, a.created_by, a.created_at,
		       p.first_name AS patient_first_name, p.last_name AS patient_last_name,
		       u.first_name AS doctor_first_name, u.last_name AS doctor_last_name, u.email AS doctor_email
		FROM appointments a
		JOIN patients p ON p.id = a.patient_id
		JOIN users u ON u.id = a.doctor_id
		WHERE a.id = $1 AND a.doctor_id = $2
	`
	if err := r.db.Get(&item, query, id, doctorID); err != nil {
		return nil, err
	}
	return &item, nil
}

type CreateAppointmentParams struct {
	PatientID int64
	DoctorID  int64
	StartAt   time.Time
	EndAt     time.Time
	Status    string
	Reason    string
	Comment   string
	CreatedBy int64
}

func (r *AppointmentRepository) Create(params CreateAppointmentParams) (*model.AppointmentWithNames, error) {
	var id int64
	query := `
		INSERT INTO appointments (patient_id, doctor_id, start_at, end_at, status, reason, comment, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	if err := r.db.Get(&id, query, params.PatientID, params.DoctorID, params.StartAt, params.EndAt, params.Status, params.Reason, params.Comment, params.CreatedBy); err != nil {
		return nil, err
	}
	return r.GetAdminByID(id)
}

type UpdateAppointmentParams struct {
	PatientID int64
	DoctorID  int64
	StartAt   time.Time
	EndAt     time.Time
	Status    string
	Reason    string
	Comment   string
}

func (r *AppointmentRepository) Update(id int64, params UpdateAppointmentParams) (*model.AppointmentWithNames, error) {
	query := `
		UPDATE appointments
		SET patient_id = $1, doctor_id = $2, start_at = $3, end_at = $4, status = $5, reason = $6, comment = $7
		WHERE id = $8
	`
	res, err := r.db.Exec(query, params.PatientID, params.DoctorID, params.StartAt, params.EndAt, params.Status, params.Reason, params.Comment, id)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("appointment not found")
	}
	return r.GetAdminByID(id)
}

func (r *AppointmentRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM appointments WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("appointment not found")
	}
	return nil
}

func (r *AppointmentRepository) SetStatus(id, doctorID int64, status string) (*model.AppointmentWithNames, error) {
	res, err := r.db.Exec(`UPDATE appointments SET status = $1 WHERE id = $2 AND doctor_id = $3`, status, id, doctorID)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("appointment not found")
	}
	return r.GetDoctorByID(id, doctorID)
}

func (r *AppointmentRepository) HasOverlap(doctorID, excludeID int64, startAt, endAt time.Time) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM appointments
		WHERE doctor_id = $1
		  AND status IN ('SCHEDULED', 'CONFIRMED')
		  AND start_at < $2
		  AND end_at > $3
	`
	args := []interface{}{doctorID, endAt, startAt}
	if excludeID > 0 {
		query += ` AND id <> $4`
		args = append(args, excludeID)
	}
	if err := r.db.Get(&count, query, args...); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AppointmentRepository) FitsSchedule(doctorID int64, startAt, endAt time.Time) (bool, error) {
	weekday := int(startAt.Weekday())
	if weekday == 0 {
		weekday = 6
	} else {
		weekday = weekday - 1
	}
	var count int
	query := `
		SELECT COUNT(*)
		FROM doctor_schedules
		WHERE doctor_id = $1
		  AND weekday = $2
		  AND start_time <= $3::time
		  AND end_time >= $4::time
	`
	if err := r.db.Get(&count, query, doctorID, weekday, startAt.Format("15:04:05"), endAt.Format("15:04:05")); err != nil {
		return false, err
	}
	return count > 0, nil
}
