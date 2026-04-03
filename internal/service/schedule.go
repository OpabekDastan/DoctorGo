package service

import (
	"errors"
	"strings"
	"doctorgo/internal/model"
	"doctorgo/internal/repository/postgres"
)

type ScheduleService struct {
	repo *postgres.ScheduleRepository
}

func NewScheduleService(repo *postgres.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) List(doctorID int64) ([]model.DoctorSchedule, error) {
	return s.repo.ListByDoctor(doctorID)
}
func (s *ScheduleService) Get(id, doctorID int64) (*model.DoctorSchedule, error) {
	return s.repo.GetByID(id, doctorID)
}

type UpsertScheduleInput struct {
	Weekday int `json:"weekday"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	SlotMinutes int `json:"slot_minutes"`
}

func normalizeScheduleInput(input *UpsertScheduleInput) {
	if input.SlotMinutes == 0 {
		input.SlotMinutes = 30
	}
}

func (s *ScheduleService) validate(input UpsertScheduleInput) error {
	if input.Weekday < 0 || input.Weekday > 6 {
		return errors.New("weekday must be between 0 and 6")
	}
	if strings.TrimSpace(input.StartTime) == "" || strings.TrimSpace(input.EndTime) == "" {
		return errors.New("start_time and end_time are required")
	}
	if input.SlotMinutes <= 0 {
		return errors.New("slot_minutes must be greater than 0")
	}
	return nil
}

func (s *ScheduleService) Create(doctorID int64, input UpsertScheduleInput) (*model.DoctorSchedule, error) {
	normalizeScheduleInput(&input)
	if err := s.validate(input); err != nil {
		return nil, err
	}
	return s.repo.Create(doctorID, postgres.UpsertScheduleParams(input))
}

func (s *ScheduleService) Update(id, doctorID int64, input UpsertScheduleInput) (*model.DoctorSchedule, error) {
	normalizeScheduleInput(&input)
	if err := s.validate(input); err != nil {
		return nil, err
	}
	return s.repo.Update(id, doctorID, postgres.UpsertScheduleParams(input))
}

func (s *ScheduleService) Delete(id, doctorID int64) error { 
	return s.repo.Delete(id, doctorID)
}
