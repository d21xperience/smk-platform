package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AcademicScore merepresentasikan satu baris nilai asesmen (DAILY/MIDTERM/FINAL)
type AcademicScore struct {
	ID                string     `json:"id"`
	Semester          int        `json:"semester"`
	AssessmentType    string     `json:"assessmentType"`
	Score             float64    `json:"score"`
	AssessmentDate    *time.Time `json:"assessmentDate,omitempty"`
	Catatan           string     `json:"catatan"`
	MataPelajaranID   int        `json:"mataPelajaranId"`
	NamaMataPelajaran string     `json:"namaMataPelajaran"`
}

// StudentBukuIndukDetail adalah DTO lengkap untuk Buku Induk Siswa
type StudentBukuIndukDetail struct {
	ID              string     `json:"id"`
	NIS             string     `json:"nis"`
	NISN            string     `json:"nisn"`
	Name            string     `json:"name"`
	BirthPlace      string     `json:"birthPlace"`
	BirthDate       *time.Time `json:"birthDate,omitempty"`
	Gender          string     `json:"gender"`
	Agama           string     `json:"agama"`
	Address         string     `json:"address"`
	Phone           string     `json:"phone"`
	DiterimaTanggal *time.Time `json:"diterimaTanggal,omitempty"`
	NamaAyah        string     `json:"namaAyah"`
	NamaIbu         string     `json:"namaIbu"`
	PekerjaanAyah   string     `json:"pekerjaanAyah"`
	PekerjaanIbu    string     `json:"pekerjaanIbu"`
	NamaWali        string     `json:"namaWali"`
	PekerjaanWali   string     `json:"pekerjaanWali"`

	// Dari tabel_siswa_pelengkap
	StatusDalamKel string `json:"statusDalamKel"`
	AnakKe         int    `json:"anakKe"`
	SekolahAsal    string `json:"sekolahAsal"`
	DiterimaKelas  string `json:"diterimaKelas"`
	AlamatOrtu     string `json:"alamatOrtu"`
	TeleponOrtu    string `json:"teleponOrtu"`
	AlamatWali     string `json:"alamatWali"`
	TeleponWali    string `json:"teleponWali"`
	FotoURL        string `json:"fotoUrl"`

	// Akademik (di-load terpisah)
	AcademicScores []AcademicScore `json:"academicScores"`

	// Kelengkapan (dihitung di service)
	SkorKelengkapan   int    `json:"skorKelengkapan"`
	StatusKelengkapan string `json:"statusKelengkapan"`
}

// StudentBukuIndukRepository adalah kontrak akses data Buku Induk Siswa
type StudentBukuIndukRepository interface {
	ListBukuInduk(ctx context.Context, schoolID, academicPeriodID string) ([]StudentBukuIndukDetail, error)
	GetBukuIndukDetail(ctx context.Context, studentID, schoolID, academicPeriodID string) (*StudentBukuIndukDetail, error)
	ListAcademicScores(ctx context.Context, studentID, schoolID, academicPeriodID string) ([]AcademicScore, error)
}

type pgxStudentBukuIndukRepository struct {
	db *pgxpool.Pool
}

func NewStudentBukuIndukRepository(db *pgxpool.Pool) StudentBukuIndukRepository {
	return &pgxStudentBukuIndukRepository{db: db}
}

// ListBukuInduk mengembalikan siswa yang punya minimal 1 nilai di school_id + academic_period_id tersebut.
// Catatan: tabel_siswa tidak punya kolom school_id (DAPODIK master), sehingga filter school_id
// dilakukan via student_academic_scores (tabel transaksi yang punya school_id).
func (r *pgxStudentBukuIndukRepository) ListBukuInduk(ctx context.Context, schoolID, academicPeriodID string) ([]StudentBukuIndukDetail, error) {
	query := `
		SELECT
			ts.peserta_didik_id::text,
			ts.nis,
			COALESCE(ts.nisn, '')      AS nisn,
			ts.nm_siswa                AS name,
			COALESCE(ts.tempat_lahir, '') AS birth_place,
			ts.tanggal_lahir,
			COALESCE(ts.jenis_kelamin, '') AS gender
		FROM tabel_siswa ts
		WHERE EXISTS (
			SELECT 1 FROM student_academic_scores sas
			WHERE sas.student_id = ts.peserta_didik_id
			  AND sas.school_id = $1
			  AND sas.academic_period_id = $2
		)
		ORDER BY ts.nm_siswa ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []StudentBukuIndukDetail
	for rows.Next() {
		var d StudentBukuIndukDetail
		if err := rows.Scan(
			&d.ID, &d.NIS, &d.NISN, &d.Name,
			&d.BirthPlace, &d.BirthDate, &d.Gender,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

// GetBukuIndukDetail mengambil 1 siswa lengkap + nilainya
func (r *pgxStudentBukuIndukRepository) GetBukuIndukDetail(ctx context.Context, studentID, schoolID, academicPeriodID string) (*StudentBukuIndukDetail, error) {
	query := `
		SELECT
			ts.peserta_didik_id::text,
			ts.nis,
			COALESCE(ts.nisn, '')             AS nisn,
			ts.nm_siswa                       AS name,
			COALESCE(ts.tempat_lahir, '')     AS birth_place,
			ts.tanggal_lahir,
			COALESCE(ts.jenis_kelamin, '')    AS gender,
			COALESCE(ts.agama, '')            AS agama,
			COALESCE(ts.alamat_siswa, '')     AS address,
			COALESCE(ts.telepon_siswa, '')    AS phone,
			ts.diterima_tanggal,
			COALESCE(ts.nm_ayah, '')          AS nm_ayah,
			COALESCE(ts.nm_ibu, '')           AS nm_ibu,
			COALESCE(ts.pekerjaan_ayah, '')   AS pekerjaan_ayah,
			COALESCE(ts.pekerjaan_ibu, '')    AS pekerjaan_ibu,
			COALESCE(ts.nm_wali, '')          AS nm_wali,
			COALESCE(ts.pekerjaan_wali, '')   AS pekerjaan_wali,
			COALESCE(tsp.status_dalam_kel, '') AS status_dalam_kel,
			COALESCE(tsp.anak_ke, 0)          AS anak_ke,
			COALESCE(tsp.sekolah_asal, '')    AS sekolah_asal,
			COALESCE(tsp.diterima_kelas, '')  AS diterima_kelas,
			COALESCE(tsp.alamat_ortu, '')     AS alamat_ortu,
			COALESCE(tsp.telepon_ortu, '')    AS telepon_ortu,
			COALESCE(tsp.alamat_wali, '')     AS alamat_wali,
			COALESCE(tsp.telepon_wali, '')    AS telepon_wali,
			COALESCE(tsp.foto_siswa, '')      AS foto_siswa
		FROM tabel_siswa ts
		LEFT JOIN tabel_siswa_pelengkap tsp ON tsp.peserta_didik_id = ts.peserta_didik_id
		WHERE ts.peserta_didik_id = $1
		LIMIT 1
	`
	var d StudentBukuIndukDetail
	err := r.db.QueryRow(ctx, query, studentID).Scan(
		&d.ID, &d.NIS, &d.NISN, &d.Name,
		&d.BirthPlace, &d.BirthDate, &d.Gender,
		&d.Agama, &d.Address, &d.Phone, &d.DiterimaTanggal,
		&d.NamaAyah, &d.NamaIbu, &d.PekerjaanAyah, &d.PekerjaanIbu,
		&d.NamaWali, &d.PekerjaanWali,
		&d.StatusDalamKel, &d.AnakKe, &d.SekolahAsal, &d.DiterimaKelas,
		&d.AlamatOrtu, &d.TeleponOrtu, &d.AlamatWali, &d.TeleponWali, &d.FotoURL,
	)
	if err != nil {
		return nil, err
	}

	// Load nilai rapor (filter school + period)
	scores, err := r.ListAcademicScores(ctx, studentID, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	d.AcademicScores = scores
	return &d, nil
}

// ListAcademicScores mengambil semua nilai untuk 1 siswa
func (r *pgxStudentBukuIndukRepository) ListAcademicScores(ctx context.Context, studentID, schoolID, academicPeriodID string) ([]AcademicScore, error) {
	query := `
		SELECT
			id::text,
			semester,
			assessment_type,
			score::float8,
			assessment_date,
			COALESCE(catatan, '') AS catatan,
			mata_pelajaran_id,
			nama_mata_pelajaran
		FROM student_academic_scores
		WHERE student_id = $1
		  AND school_id = $2
		  AND academic_period_id = $3
		ORDER BY semester ASC, nama_mata_pelajaran ASC, assessment_type ASC
	`
	rows, err := r.db.Query(ctx, query, studentID, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []AcademicScore
	for rows.Next() {
		var s AcademicScore
		if err := rows.Scan(
			&s.ID, &s.Semester, &s.AssessmentType, &s.Score,
			&s.AssessmentDate, &s.Catatan, &s.MataPelajaranID, &s.NamaMataPelajaran,
		); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, nil
}
