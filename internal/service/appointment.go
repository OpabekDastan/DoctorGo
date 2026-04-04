package service

import (
	"errors"
	"time"

	"doctor_go/internal/model"
	"doctor_go/internal/repository/postgres"
)

type AppointmentService struct {
	repo *postgres.AppointmentRepository
}

func NewAppointmentService(repo *postgres.AppointmentRepository) *AppointmentService {
	return &AppointmentService{repo: repo}
}

func (s *AppointmentService) ListAdmin() ([]model.AppointmentWithNames, error) {
	return s.repo.ListAdmin()
}
func (s *AppointmentService) ListDoctor(doctorID int64) ([]model.AppointmentWithNames, error) {
	return s.repo.ListByDoctor(doctorID)
}
func (s *AppointmentService) GetAdmin(id int64) (*model.AppointmentWithNames, error) {
	return s.repo.GetAdminByID(id)
}
func (s *AppointmentService) GetDoctor(id, doctorID int64) (*model.AppointmentWithNames, error) {
	return s.repo.GetDoctorByID(id, doctorID)
}

type UpsertAppointmentInput struct {
	PatientID int64     `json:"patient_id"`
	DoctorID  int64     `json:"doctor_id"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	Comment   string    `json:"comment"`
}

func (s *AppointmentService) validate(excludeID int64, input UpsertAppointmentInput) error {
	if input.PatientID == 0 || input.DoctorID == 0 {
		return errors.New("patient_id and doctor_id are required")
	}
	if !input.EndAt.After(input.StartAt) {
		return errors.New("end_at must be after start_at")
	}
	if input.Status == "" {
		input.Status = model.AppointmentScheduled
	}
	overlap, err := s.repo.HasOverlap(input.DoctorID, excludeID, input.StartAt, input.EndAt)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("appointment overlaps with another active appointment for this doctor")
	}
	fits, err := s.repo.FitsSchedule(input.DoctorID, input.StartAt, input.EndAt)
	if err != nil {
		return err
	}
	if !fits {
		return errors.New("appointment is outside doctor's working hours")
	}
	return nil
}

func (s *AppointmentService) Create(createdBy int64, input UpsertAppointmentInput) (*model.AppointmentWithNames, error) {
	if input.Status == "" {
		input.Status = model.AppointmentScheduled
	}
	if err := s.validate(0, input); err != nil {
		return nil, err
	}
	return s.repo.Create(postgres.CreateAppointmentParams{
		PatientID: input.PatientID,
		DoctorID:  input.DoctorID,
		StartAt:   input.StartAt,
		EndAt:     input.EndAt,
		Status:    input.Status,
		Reason:    input.Reason,
		Comment:   input.Comment,
		CreatedBy: createdBy,
	})
}

func (s *AppointmentService) Update(id int64, input UpsertAppointmentInput) (*model.AppointmentWithNames, error) {
	if input.Status == "" {
		input.Status = model.AppointmentScheduled
	}
	if err := s.validate(id, input); err != nil {
		return nil, err
	}
	return s.repo.Update(id, postgres.UpdateAppointmentParams{
		PatientID: input.PatientID,
		DoctorID:  input.DoctorID,
		StartAt:   input.StartAt,
		EndAt:     input.EndAt,
		Status:    input.Status,
		Reason:    input.Reason,
		Comment:   input.Comment,
	})
}

func (s *AppointmentService) Delete(id int64) error { return s.repo.Delete(id) }

func (s *AppointmentService) SetStatus(id, doctorID int64, status string) (*model.AppointmentWithNames, error) {
	allowed := map[string]bool{
		model.AppointmentScheduled: true,
		model.AppointmentConfirmed: true,
		model.AppointmentCompleted: true,
		model.AppointmentCancelled: true,
		model.AppointmentNoShow:    true,
	}
	if !allowed[status] {
		return nil, errors.New("invalid status")
	}
	return s.repo.SetStatus(id, doctorID, status)
}
