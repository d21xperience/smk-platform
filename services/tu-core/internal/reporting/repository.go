package reporting

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ReportFilter membawa konteks operasional wajib (Prinsip 5).
type ReportFilter struct {
	SchoolID         string
	AcademicPeriodID string
	StartDate        string // Format: YYYY-MM-DD
	EndDate          string // Format: YYYY-MM-DD
	ClassID          string // Opsional: untuk filter per kelas
}

// ReportingRepository adalah kontrak untuk mengakses Read Model (Projection Tables).
type ReportingRepository interface {
	GetAttendanceSummary(ctx context.Context, filter ReportFilter) ([]map[string]interface{}, error)
	GetFinanceSummary(ctx context.Context, filter ReportFilter) (map[string]interface{}, error)
}

// pgxReportingRepository mengimplementasikan ReportingRepository menggunakan pgxpool.
type pgxReportingRepository struct {
	db *pgxpool.Pool
}

// NewReportingRepository membuat instance baru dari repository.
func NewReportingRepository(db *pgxpool.Pool) ReportingRepository {
	return &pgxReportingRepository{db: db}
}

// GetAttendanceSummary mengambil data agregat absensi harian dari read model.
func (r *pgxReportingRepository) GetAttendanceSummary(ctx context.Context, filter ReportFilter) ([]map[string]interface{}, error) {
	// Validasi Operational Context (Prinsip 5)
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, fmt.Errorf("schoolId dan academicPeriodId wajib diisi")
	}

	query := `
		SELECT
			date,
			SUM(total_students) as total_students,
			SUM(present_count) as present_count,
			SUM(absent_count) as absent_count,
			SUM(late_count) as late_count,
			ROUND(AVG(attendance_rate), 2) as avg_attendance_rate
		FROM reporting_attendance_daily
		WHERE school_id = $1
		  AND academic_period_id = $2
		  AND date BETWEEN $3 AND $4
	`
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID, filter.StartDate, filter.EndDate}

	if filter.ClassID != "" {
		query += " AND class_id = $5"
		args = append(args, filter.ClassID)
	}

	query += " GROUP BY date ORDER BY date ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query attendance summary: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var date string
		var total, present, absent, late int
		var avgRate float64

		if err := rows.Scan(&date, &total, &present, &absent, &late, &avgRate); err != nil {
			return nil, fmt.Errorf("gagal scan row attendance summary: %w", err)
		}

		results = append(results, map[string]interface{}{
			"date":                date,
			"total_students":      total,
			"present_count":       present,
			"absent_count":        absent,
			"late_count":          late,
			"avg_attendance_rate": avgRate,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance summary rows: %w", err)
	}

	return results, nil
}

// GetFinanceSummary mengambil data agregat keuangan bulanan dari read model.
func (r *pgxReportingRepository) GetFinanceSummary(ctx context.Context, filter ReportFilter) (map[string]interface{}, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, fmt.Errorf("schoolId dan academicPeriodId wajib diisi")
	}

	query := `
		SELECT
			component_type,
			SUM(total_billed) as total_billed,
			SUM(total_paid) as total_paid,
			SUM(total_outstanding) as total_outstanding,
			SUM(payment_count) as total_payments
		FROM reporting_finance_monthly
		WHERE school_id = $1
		  AND academic_period_id = $2
		GROUP BY component_type
	`

	rows, err := r.db.Query(ctx, query, filter.SchoolID, filter.AcademicPeriodID)
	if err != nil {
		return nil, fmt.Errorf("gagal query finance summary: %w", err)
	}
	defer rows.Close()

	summary := make(map[string]interface{})
	var components []map[string]interface{}

	for rows.Next() {
		var compType string
		var billed, paid, outstanding float64
		var count int

		if err := rows.Scan(&compType, &billed, &paid, &outstanding, &count); err != nil {
			return nil, fmt.Errorf("gagal scan row finance summary: %w", err)
		}

		components = append(components, map[string]interface{}{
			"component_type":    compType,
			"total_billed":      billed,
			"total_paid":        paid,
			"total_outstanding": outstanding,
			"payment_count":     count,
		})
	}

	summary["components"] = components
	return summary, rows.Err()
}
