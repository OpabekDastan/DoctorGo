package repository

import "github.com/jmoiron/sqlx"

type DashboardRepository struct {
	db *sqlx.DB
}

func NewDashboardRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

type GeneralStats struct {
	TotalDoctors      int `db:"total_doctors"      json:"total_doctors"`
	TotalPatients     int `db:"total_patients"     json:"total_patients"`
	TotalAppointments int `db:"total_appointments" json:"total_appointments"`
}

// GetGeneralStats возвращает общую статистику по клинике одним запросом.
func (r *DashboardRepository) GetGeneralStats() (*GeneralStats, error) {
	var stats GeneralStats
	query := `
		SELECT
			(SELECT COUNT(*) FROM users    WHERE role = 'DOCTOR') AS total_doctors,
			(SELECT COUNT(*) FROM patients)                        AS total_patients,
			(SELECT COUNT(*) FROM appointments)                    AS total_appointments
	`
	if err := r.db.Get(&stats, query); err != nil {
		return nil, err
	}
	return &stats, nil
}
