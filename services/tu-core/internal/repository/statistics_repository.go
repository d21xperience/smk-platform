package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GradeStat struct {
	Tingkat   int `json:"tingkat"`
	Total     int `json:"total"`
	LakiLaki  int `json:"laki_laki"`
	Perempuan int `json:"perempuan"`
}

type MajorStat struct {
	Jurusan string `json:"jurusan"`
	Total   int    `json:"total"`
}

type StatisticsData struct {
	TotalSiswa      int         `json:"total_siswa"`
	TotalLakiLaki   int         `json:"total_laki_laki"`
	TotalPerempuan  int         `json:"total_perempuan"`
	TotalGuru       int         `json:"total_guru"`
	TotalRombel     int         `json:"total_rombel"`
	TotalJurusan    int         `json:"total_jurusan"`
	PerAngkatan     []GradeStat `json:"per_angkatan"`
	DistribusiJurusan []MajorStat `json:"distribusi_jurusan"`
}

type StatisticsRepository interface {
	GetStatistics(ctx context.Context, schoolID, periodID string) (*StatisticsData, error)
}

type pgxStatisticsRepository struct {
	db *pgxpool.Pool
}

func NewStatisticsRepository(db *pgxpool.Pool) StatisticsRepository {
	return &pgxStatisticsRepository{db: db}
}

func (r *pgxStatisticsRepository) GetStatistics(ctx context.Context, schoolID, periodID string) (*StatisticsData, error) {
	stats := &StatisticsData{
		PerAngkatan:       []GradeStat{},
		DistribusiJurusan: []MajorStat{},
	}

	// 1. Hitung Total Siswa, L/P di Semester ini
	queryTotal := `
		SELECT 
			COUNT(ts.peserta_didik_id),
			COUNT(CASE WHEN ts.jenis_kelamin IN ('L', 'MALE') THEN 1 END),
			COUNT(CASE WHEN ts.jenis_kelamin IN ('P', 'FEMALE') THEN 1 END)
		FROM tabel_anggotakelas tak
		JOIN tabel_siswa ts ON tak.peserta_didik_id = ts.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE tk.sekolah_id::text = $1 AND tak.semester_id = $2
	`
	err := r.db.QueryRow(ctx, queryTotal, schoolID, periodID).Scan(
		&stats.TotalSiswa, &stats.TotalLakiLaki, &stats.TotalPerempuan,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: gagal total statistik: %v", ErrDatabaseError, err)
	}

	// 2. Hitung Total Guru (PTK aktif)
	queryGuru := `SELECT COUNT(*) FROM tabel_ptk WHERE soft_delete = 0`
	_ = r.db.QueryRow(ctx, queryGuru).Scan(&stats.TotalGuru)

	// 3. Hitung Total Rombel & Jurusan di semester ini
	queryRombel := `
		SELECT 
			COUNT(rombongan_belajar_id), 
			COUNT(DISTINCT nama_jurusan_sp) 
		FROM tabel_kelas 
		WHERE sekolah_id::text = $1 AND semester_id = $2
	`
	_ = r.db.QueryRow(ctx, queryRombel, schoolID, periodID).Scan(&stats.TotalRombel, &stats.TotalJurusan)

	// 4. Hitung Per Angkatan (Tingkat X, XI, XII)
	queryGrade := `
		SELECT 
			tk.tingkat_pendidikan_id,
			COUNT(ts.peserta_didik_id),
			COUNT(CASE WHEN ts.jenis_kelamin IN ('L', 'MALE') THEN 1 END),
			COUNT(CASE WHEN ts.jenis_kelamin IN ('P', 'FEMALE') THEN 1 END)
		FROM tabel_anggotakelas tak
		JOIN tabel_siswa ts ON tak.peserta_didik_id = ts.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE tk.sekolah_id::text = $1 AND tak.semester_id = $2
		GROUP BY tk.tingkat_pendidikan_id
		ORDER BY tk.tingkat_pendidikan_id ASC
	`
	rows, err := r.db.Query(ctx, queryGrade, schoolID, periodID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var g GradeStat
			if err := rows.Scan(&g.Tingkat, &g.Total, &g.LakiLaki, &g.Perempuan); err == nil {
				stats.PerAngkatan = append(stats.PerAngkatan, g)
			}
		}
	}

	// 5. Distribusi Siswa Per Jurusan (Untuk Grafik Bar)
	queryJurusan := `
		SELECT 
			COALESCE(tk.nama_jurusan_sp, 'Lainnya'), 
			COUNT(ts.peserta_didik_id)
		FROM tabel_anggotakelas tak
		JOIN tabel_siswa ts ON tak.peserta_didik_id = ts.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE tk.sekolah_id::text = $1 AND tak.semester_id = $2
		GROUP BY tk.nama_jurusan_sp
		ORDER BY COUNT(ts.peserta_didik_id) DESC
	`
	rowsJur, err := r.db.Query(ctx, queryJurusan, schoolID, periodID)
	if err == nil {
		defer rowsJur.Close()
		for rowsJur.Next() {
			var m MajorStat
			if err := rowsJur.Scan(&m.Jurusan, &m.Total); err == nil {
				stats.DistribusiJurusan = append(stats.DistribusiJurusan, m)
			}
		}
	}

	return stats, nil
}
