package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JurusanStat struct {
	Jurusan      string `json:"jurusan"`
	SiswaAkhir   int    `json:"akhir"`
	MutasiMasuk  int    `json:"masuk"`
	MutasiKeluar int    `json:"keluar"`
	SiswaAwal    int    `json:"awal"`
}

type KesiswaanReportRepository interface {
	GetStatistikJurusan(ctx context.Context, schoolID, academicPeriodID string) ([]JurusanStat, error)
	GetTotalMutasi(ctx context.Context, schoolID, academicPeriodID string) (int, int, error)
}

type pgxKesiswaanReportRepository struct {
	db *pgxpool.Pool
}

func NewKesiswaanReportRepository(db *pgxpool.Pool) KesiswaanReportRepository {
	return &pgxKesiswaanReportRepository{db: db}
}

func (r *pgxKesiswaanReportRepository) GetStatistikJurusan(ctx context.Context, schoolID, academicPeriodID string) ([]JurusanStat, error) {
	// Menghitung total siswa aktif per jurusan pada semester berjalan.
	// Menggunakan fallback tj.nama_jurusan -> tk.nama_jurusan_sp -> 'Tidak Diketahui'
	query := `
		SELECT
			COALESCE(tj.nama_jurusan, tk.nama_jurusan_sp, 'Tidak Diketahui') AS nama_jurusan,
			COUNT(DISTINCT tak.peserta_didik_id) AS total_siswa
		FROM tabel_anggotakelas tak
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		LEFT JOIN tabel_jurusan tj ON tk.jurusan_id = tj.jurusan_id
		WHERE tak.semester_id = $1 AND tk.sekolah_id = $2
		GROUP BY COALESCE(tj.nama_jurusan, tk.nama_jurusan_sp, 'Tidak Diketahui')
		ORDER BY total_siswa DESC
	`
	rows, err := r.db.Query(ctx, query, academicPeriodID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []JurusanStat
	for rows.Next() {
		var s JurusanStat
		if err := rows.Scan(&s.Jurusan, &s.SiswaAkhir); err != nil {
			return nil, err
		}
		s.SiswaAwal = s.SiswaAkhir // Asumsi tanpa data historis mutasi per jurusan
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *pgxKesiswaanReportRepository) GetTotalMutasi(ctx context.Context, schoolID, academicPeriodID string) (int, int, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN mutation_type = 'MASUK' THEN 1 ELSE 0 END), 0) AS masuk,
			COALESCE(SUM(CASE WHEN mutation_type = 'KELUAR' THEN 1 ELSE 0 END), 0) AS keluar
		FROM student_mutations
		WHERE school_id = $1 AND academic_period_id = $2
	`
	var masuk, keluar int
	err := r.db.QueryRow(ctx, query, schoolID, academicPeriodID).Scan(&masuk, &keluar)
	if err != nil {
		return 0, 0, err
	}
	return masuk, keluar, nil
}
