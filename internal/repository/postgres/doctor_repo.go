package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"doctorgo/internal/model"
)

type DoctorRepository struct {
	db *sqlx.DB
}

func NewDoctorRepository(db *sqlx.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

func (r *DoctorRepository) List() ([]model.Doctor, error) {
	items := []model.Doctor{}
	query := `SELECT u.id, u.email, u.role, u.first_name, u.last_name, u.is_active, u.created_at,
		       COALESCE(dp.specialization, '') AS specialization, COALESCE(dp.phone, '') AS phone
			   FROM users u LEFT JOIN doctor_profiles dp ON dp.user_id = u.id WHERE u.role = 'DOCTOR' ORDER BY u.id`
	if err := r.db.Select(&items, query); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DoctorRepository) GetByID(id int64) (*model.Doctor, error) {
	var item model.Doctor
	query := `SELECT u.id, u.email, u.role, u.first_name, u.last_name, u.is_active, u.created_at,
		       COALESCE(dp.specialization, '') AS specialization, COALESCE(dp.phone, '') AS phone
			   FROM users u LEFT JOIN doctor_profiles dp ON dp.user_id = u.id WHERE u.role = 'DOCTOR' AND u.id = $1`
	if err := r.db.Get(&item, query, id); err != nil {
		return nil, err
	}
	return &item, nil
}

type CreateDoctorParams struct {
	Email string
	PasswordHash string
	FirstName string
	LastName string
	Specialization string
	Phone string
}

func (r *DoctorRepository) Create(params CreateDoctorParams) (*model.Doctor, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var userID int64
	userQuery := `INSERT INTO users (email, password_hash, role, first_name, last_name, is_active)
				  VALUES ($1, $2, 'DOCTOR', $3, $4, true) RETURNING id`
	if err := tx.Get(&userID, userQuery, params.Email, params.PasswordHash, params.FirstName, params.LastName); err != nil {
		return nil, err
	}

	profileQuery := `INSERT INTO doctor_profiles (user_id, specialization, phone) VALUES ($1, $2, $3)`
	if _, err := tx.Exec(profileQuery, userID, params.Specialization, params.Phone); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(userID)
}

type UpdateDoctorParams struct {
	FirstName string
	LastName string
	IsActive bool
	Specialization string
	Phone string
}

func (r *DoctorRepository) Update(id int64, params UpdateDoctorParams) (*model.Doctor, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userQuery := `UPDATE users SET first_name = $1, last_name = $2, is_active = $3 WHERE id = $4 AND role = 'DOCTOR'`
	res, err := tx.Exec(userQuery, params.FirstName, params.LastName, params.IsActive, id)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("doctor not found")
	}

	profileQuery := `UPDATE doctor_profiles SET specialization = $1, phone = $2 WHERE user_id = $3`
	if _, err := tx.Exec(profileQuery, params.Specialization, params.Phone, id); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *DoctorRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = $1 AND role = 'DOCTOR'`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("doctor not found")
	}
	return nil
}
