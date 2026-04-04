package service

import (
	"errors"

	"doctor_go/internal/model"
	"doctor_go/internal/repository/postgres"
)

type PatientService struct {
	repo *postgres.PatientRepository
}

func NewPatientService(repo *postgres.PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) List() ([]model.Patient, error)       { return s.repo.List() }
func (s *PatientService) Get(id int64) (*model.Patient, error) { return s.repo.GetByID(id) }

type UpsertPatientInput struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	BirthDate *string `json:"birth_date"`
	Gender    string  `json:"gender"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	Address   string  `json:"address"`
	Comment   string  `json:"comment"`
}

func (s *PatientService) Create(input UpsertPatientInput) (*model.Patient, error) {
	if input.FirstName == "" || input.LastName == "" {
		return nil, errors.New("first_name and last_name are required")
	}
	if input.Gender == "" {
		input.Gender = "U"
	}
	return s.repo.Create(postgres.CreatePatientParams(input))
}

func (s *PatientService) Update(id int64, input UpsertPatientInput) (*model.Patient, error) {
	if input.FirstName == "" || input.LastName == "" {
		return nil, errors.New("first_name and last_name are required")
	}
	if input.Gender == "" {
		input.Gender = "U"
	}
	return s.repo.Update(id, postgres.UpdatePatientParams(input))
}

func (s *PatientService) Delete(id int64) error { return s.repo.Delete(id) }
