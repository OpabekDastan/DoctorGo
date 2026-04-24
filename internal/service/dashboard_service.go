package service

import "doctor_go/internal/repository/postgres"

type DashboardService struct {
	repo *postgres.DashboardRepository
}

func NewDashboardService(repo *postgres.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetStats() (*postgres.GeneralStats, error) {
	return s.repo.GetGeneralStats()
}
