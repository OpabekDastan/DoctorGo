package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"doctor_go/internal/model"
)

type ScheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) ListByDoctor(doctorID int64) ([]model.DoctorSchedule, error) {
	items := []model.DoctorSchedule{}
	query := `SELECT id, doctor_id, weekday, start_time::text AS start_time, end_time::text AS end_time, slot_minutes
			  FROM doctor_schedules WHERE doctor_id = $1 ORDER BY weekday, start_time`
	if err := r.db.Select(&items, query, doctorID); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ScheduleRepository) GetByID(id, doctorID int64) (*model.DoctorSchedule, error) {
	var item model.DoctorSchedule
	query := `SELECT id, doctor_id, weekday, start_time::text AS start_time, end_time::text AS end_time, slot_minutes
			  FROM doctor_schedules WHERE id = $1 AND doctor_id = $2`
	if err := r.db.Get(&item, query, id, doctorID); err != nil {
		return nil, err
	}
	return &item, nil
}

type UpsertScheduleParams struct {
	Weekday int
	StartTime string
	EndTime string
	SlotMinutes int
}

func (r *ScheduleRepository) Create(doctorID int64, params UpsertScheduleParams) (*model.DoctorSchedule, error) {
	var id int64
	query := `INSERT INTO doctor_schedules (doctor_id, weekday, start_time, end_time, slot_minutes) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	if err := r.db.Get(&id, query, doctorID, params.Weekday, params.StartTime, params.EndTime, params.SlotMinutes); err != nil {
		return nil, err
	}
	return r.GetByID(id, doctorID)
}

func (r *ScheduleRepository) Update(id, doctorID int64, params UpsertScheduleParams) (*model.DoctorSchedule, error) {
	query := `UPDATE doctor_schedules SET weekday = $1, start_time = $2, end_time = $3, slot_minutes = $4 WHERE id = $5 AND doctor_id = $6`
	res, err := r.db.Exec(query, params.Weekday, params.StartTime, params.EndTime, params.SlotMinutes, id, doctorID)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("schedule not found")
	}
	return r.GetByID(id, doctorID)
}

func (r *ScheduleRepository) Delete(id, doctorID int64) error {
	res, err := r.db.Exec(`DELETE FROM doctor_schedules WHERE id = $1 AND doctor_id = $2`, id, doctorID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("schedule not found")
	}
	return nil
}
