package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherWorkload struct {
	TeacherName string `json:"teacherName"`
	SubjectName string `json:"subjectName"`
	JP          int    `json:"jp"` // Menggunakan jumlah kelas yg diajar sebagai proxy "Jam Pelajaran"
}

type KurikulumStats struct {
	TotalMapel        int               `json:"totalMapel"`
	GuruMengajarAktif int               `json:"guruMengajarAktif"`
	BebanMengajar     []TeacherWorkload `json:"bebanMengajar"`
}

type KurikulumStatRepository interface {
	GetStats(ctx context.Context, schoolID, periodID string) (*KurikulumStats, error)
}

type pgxKurikulumStatRepository struct {
	db *pgxpool.Pool
}

func NewKurikulumStatRepository(db *pgxpool.Pool) KurikulumStatRepository {
	return &pgxKurikulumStatRepository{db: db}
}

func (r *pgxKurikulumStatRepository) GetStats(ctx context.Context, schoolID, periodID string) (*KurikulumStats, error) {
	stats := &KurikulumStats{BebanMengajar: []TeacherWorkload{}}

	// 1. Total Mapel (Global)
	_ = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM tabel_mapel").Scan(&stats.TotalMapel)

	// 2. Guru Mengajar Aktif (Distinct PTK di tabel_pembelajaran)
	queryGuru := `
		SELECT COUNT(DISTINCT tpm.ptk_terdaftar_id)
		FROM tabel_pembelajaran tpm
		JOIN tabel_kelas tk ON tpm.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE tk.sekolah_id::text = $1 AND tpm.semester_id = $2 AND tpm.ptk_terdaftar_id IS NOT NULL
	`
	_ = r.db.QueryRow(ctx, queryGuru, schoolID, periodID).Scan(&stats.GuruMengajarAktif)

	// 3. Beban Mengajar (Dihitung dari jumlah kelas yang diajar per guru)
	queryBeban := `
		SELECT tp.nama, COUNT(tpm.pembelajaran_id) as total_jp
		FROM tabel_pembelajaran tpm
		JOIN tabel_kelas tk ON tpm.rombongan_belajar_id = tk.rombongan_belajar_id
		JOIN tabel_ptk tp ON tpm.ptk_terdaftar_id = tp.ptk_id
		WHERE tk.sekolah_id::text = $1 AND tpm.semester_id = $2
		GROUP BY tp.nama
		ORDER BY total_jp DESC
		LIMIT 5
	`
	rows, err := r.db.Query(ctx, queryBeban, schoolID, periodID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w TeacherWorkload
			if err := rows.Scan(&w.TeacherName, &w.JP); err == nil {
				w.SubjectName = "Multi Mapel" // Karena kita agregasi semua mapel yg diajar
				stats.BebanMengajar = append(stats.BebanMengajar, w)
			}
		}
	}

	return stats, nil
}
