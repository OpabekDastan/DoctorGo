package service

import (
	"errors"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"doctor_go/internal/model"
	"doctor_go/internal/repository/postgres"
)

type DoctorService struct {
	repo *postgres.DoctorRepository
}

func NewDoctorService(repo *postgres.DoctorRepository) *DoctorService {
	return &DoctorService{repo: repo}
}

func (s *DoctorService) List() ([]model.Doctor, error){
	return s.repo.List() 
}
func (s *DoctorService) Get(id int64) (*model.Doctor, error) {
	return s.repo.GetByID(id) 
}

type CreateDoctorInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Specialization string `json:"specialization"`
	Phone string `json:"phone"`
}

func (s *DoctorService) Create(input CreateDoctorInput) (*model.Doctor, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Password = strings.TrimSpace(input.Password)
	if input.Email == "" || input.Password == "" {
		return nil, errors.New("email and password are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(postgres.CreateDoctorParams{
		Email: input.Email,
		PasswordHash: string(hash),
		FirstName: input.FirstName,
		LastName: input.LastName,
		Specialization: input.Specialization,
		Phone: input.Phone,
	})
}

type UpdateDoctorInput struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	IsActive bool   `json:"is_active"`
	Specialization string `json:"specialization"`
	Phone string `json:"phone"`
}

func (s *DoctorService) Update(id int64, input UpdateDoctorInput) (*model.Doctor, error) {
	return s.repo.Update(id, postgres.UpdateDoctorParams{
		FirstName: input.FirstName,
		LastName: input.LastName,
		IsActive: input.IsActive,
		Specialization: input.Specialization,
		Phone:input.Phone,
	})
}

func (s *DoctorService) Delete(id int64) error { 
	return s.repo.Delete(id) 
}
