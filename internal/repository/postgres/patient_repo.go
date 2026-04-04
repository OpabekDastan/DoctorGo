package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"doctor_go/internal/model"
)

type PatientRepository struct {
	db *sqlx.DB
}

func NewPatientRepository(db *sqlx.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) List() ([]model.Patient, error) {
	items := []model.Patient{}
	query := `
		SELECT id, first_name, last_name, birth_date::text AS birth_date, gender, phone, email, address, comment, created_at
		FROM patients
		ORDER BY id DESC
	`
	if err := r.db.Select(&items, query); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PatientRepository) GetByID(id int64) (*model.Patient, error) {
	var item model.Patient
	query := `
		SELECT id, first_name, last_name, birth_date::text AS birth_date, gender, phone, email, address, comment, created_at
		FROM patients
		WHERE id = $1
	`
	if err := r.db.Get(&item, query, id); err != nil {
		return nil, err
	}
	return &item, nil
}

type CreatePatientParams struct {
	FirstName string
	LastName  string
	BirthDate *string
	Gender    string
	Phone     string
	Email     string
	Address   string
	Comment   string
}

func (r *PatientRepository) Create(params CreatePatientParams) (*model.Patient, error) {
	var id int64
	query := `
		INSERT INTO patients (first_name, last_name, birth_date, gender, phone, email, address, comment)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	if err := r.db.Get(&id, query, params.FirstName, params.LastName, params.BirthDate, params.Gender, params.Phone, params.Email, params.Address, params.Comment); err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

type UpdatePatientParams = CreatePatientParams

func (r *PatientRepository) Update(id int64, params UpdatePatientParams) (*model.Patient, error) {
	query := `
		UPDATE patients
		SET first_name = $1, last_name = $2, birth_date = $3, gender = $4, phone = $5, email = $6, address = $7, comment = $8
		WHERE id = $9
	`
	res, err := r.db.Exec(query, params.FirstName, params.LastName, params.BirthDate, params.Gender, params.Phone, params.Email, params.Address, params.Comment, id)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("patient not found")
	}
	return r.GetByID(id)
}

func (r *PatientRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM patients WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("patient not found")
	}
	return nil
}
