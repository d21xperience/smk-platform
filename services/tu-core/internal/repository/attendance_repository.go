package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceSummary struct {
	TotalSessions        int
	PresentCount         int
	AbsentCount          int
	LateCount            int
	SickCount            int
	PermissionCount      int
	AttendancePercentage float64
}

type AttendanceRepository interface {
	GetStudentAttendanceSummary(ctx context.Context, schoolID, academicPeriodID, studentID string) (*AttendanceSummary, error)
}

type pgxAttendanceRepository struct {
	db *pgxpool.Pool
}

func NewAttendanceRepository(db *pgxpool.Pool) AttendanceRepository {
	return &pgxAttendanceRepository{db: db}
}

func (r *pgxAttendanceRepository) GetStudentAttendanceSummary(ctx context.Context, schoolID, academicPeriodID, studentID string) (*AttendanceSummary, error) {
	// Query ke tabel_kehadiran DAPODIK
	query := `
		SELECT
			COALESCE(SUM(tk.sakit), 0) as sakit,
			COALESCE(SUM(tk.izin), 0) as izin,
			COALESCE(SUM(tk.tanpa_keterangan), 0) as alpa
		FROM tabel_kehadiran tk
		JOIN tabel_anggotakelas tak ON tk.anggota_rombel_id = tak.anggota_rombel_id
		WHERE tak.peserta_didik_id = $1
		  AND tak.semester_id = $2
		  AND tak.rombongan_belajar_id IN (
			  SELECT rombongan_belajar_id FROM tabel_kelas WHERE sekolah_id = $3
		  )
	`

	var sakit, izin, alpa int
	err := r.db.QueryRow(ctx, query, studentID, academicPeriodID, schoolID).Scan(&sakit, &izin, &alpa)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &AttendanceSummary{}, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Karena DAPODIK tidak menyimpan "Hadir" dan "Terlambat" secara eksplisit,
	// kita asumsikan TotalSessions adalah jumlah dari Sakit + Izin + Alpa untuk sementara.
	// (Atau kita bisa biarkan 0 untuk Hadir/Terlambat)
	total := sakit + izin + alpa
	summary := &AttendanceSummary{
		TotalSessions:   total,
		SickCount:       sakit,
		PermissionCount: izin,
		AbsentCount:     alpa,
		PresentCount:    0, // Belum ada data hadir di DAPODIK
		LateCount:       0, // Belum ada data terlambat di DAPODIK
	}

	if summary.TotalSessions > 0 {
		summary.AttendancePercentage = float64(summary.PresentCount+summary.LateCount) / float64(summary.TotalSessions) * 100
	}

	return summary, nil
}
