package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository interface {
	GetKPIStats(ctx context.Context, schoolID string, semesterID string) (map[string]interface{}, error)
	GetJurusanDistribusi(ctx context.Context, schoolID string, semesterID string) ([]map[string]interface{}, error)
	GetKehadiranHariIni(ctx context.Context, schoolID string, semesterID string) ([]map[string]interface{}, error)
	GetAgendaMendatang(ctx context.Context, schoolID string, semesterID string) ([]map[string]interface{}, error)
}

type dashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetKPIStats(ctx context.Context, schoolID string, semesterID string) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"totalSiswa":  0,
		"totalGuru":   0,
		"totalRombel": 0,
	}

	var totalSiswa int
	err := r.db.QueryRow(ctx, `SELECT COUNT(DISTINCT siswa_id) FROM tabel_anggotakelas WHERE semester_id = $1`, semesterID).Scan(&totalSiswa)
	if err == nil { stats["totalSiswa"] = totalSiswa }

	var totalRombel int
	err = r.db.QueryRow(ctx, `SELECT COUNT(DISTINCT kelas_id) FROM tabel_anggotakelas WHERE semester_id = $1`, semesterID).Scan(&totalRombel)
	if err == nil { stats["totalRombel"] = totalRombel }

	var totalGuru int
	err = r.db.QueryRow(ctx, `SELECT COUNT(ptk_id) FROM tabel_ptk WHERE status_kepegawaian = 'Aktif'`).Scan(&totalGuru)
	if err == nil { stats["totalGuru"] = totalGuru }

	return stats, nil
}

func (r *dashboardRepository) GetJurusanDistribusi(ctx context.Context, schoolID string, semesterID string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT tk.nama_jurusan, COUNT(tak.siswa_id) as jumlah
		FROM tabel_anggotakelas tak
		JOIN tabel_kelas tk ON tak.kelas_id = tk.rombongan_belajar_id
		WHERE tak.semester_id = $1
		GROUP BY tk.nama_jurusan ORDER BY jumlah DESC`, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var nama string
		var jumlah int
		if err := rows.Scan(&nama, &jumlah); err == nil {
			result = append(result, map[string]interface{}{"nama": nama, "jumlah": jumlah, "warna": "#0a192f"})
		}
	}
	return result, nil
}

func (r *dashboardRepository) GetKehadiranHariIni(ctx context.Context, schoolID string, semesterID string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ar.status, COUNT(ar.id) as jumlah
		FROM attendance_records ar
		JOIN attendance_sessions asess ON ar.session_id = asess.id
		WHERE asess.academic_period_id = $1 AND asess.session_date = CURRENT_DATE
		GROUP BY ar.status`, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	counts := map[string]int{"PRESENT": 0, "PERMISSION": 0, "SICK": 0, "ABSENT": 0}
	total := 0
	for rows.Next() {
		var status string
		var jumlah int
		if err := rows.Scan(&status, &jumlah); err == nil {
			counts[status] = jumlah
			total += jumlah
		}
	}
	if total == 0 { total = 1 }

	return []map[string]interface{}{
		{"label": "Hadir", "persen": (counts["PRESENT"] * 100) / total, "jumlah": counts["PRESENT"], "warna": "#2e7d32"},
		{"label": "Izin", "persen": (counts["PERMISSION"] * 100) / total, "jumlah": counts["PERMISSION"], "warna": "#1976d2"},
		{"label": "Sakit", "persen": (counts["SICK"] * 100) / total, "jumlah": counts["SICK"], "warna": "#f9a825"},
		{"label": "Alpa", "persen": (counts["ABSENT"] * 100) / total, "jumlah": counts["ABSENT"], "warna": "#c62828"},
	}, nil
}

func (r *dashboardRepository) GetAgendaMendatang(ctx context.Context, schoolID string, semesterID string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT jenis, nomor_agenda, tanggal_catat 
		FROM agenda_entries WHERE school_id = $1 ORDER BY tanggal_catat DESC LIMIT 3`, schoolID)
	if err != nil { return nil, err }
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var jenis, nomor string
		var tanggal interface{}
		if err := rows.Scan(&jenis, &nomor, &tanggal); err == nil {
			result = append(result, map[string]interface{}{
				"judul": "Agenda " + jenis + " - " + nomor, "waktu": "Hari ini", "tanggal": "14", "bulan": "Jul",
			})
		}
	}
	return result, nil
}
