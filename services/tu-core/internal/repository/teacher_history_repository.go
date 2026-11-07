package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EducationHistory struct {
	ID          string `json:"id"`
	PtkID       string `json:"ptkId"`
	Jenjang     string `json:"jenjang"`
	Jurusan     string `json:"jurusan"`
	NamaSekolah string `json:"namaSekolah"`
	TahunLulus  int    `json:"tahunLulus"`
	Gelar       string `json:"gelar"`
	DokumenURL  string `json:"dokumenUrl"`
}

type CertificationHistory struct {
	ID               string `json:"id"`
	PtkID            string `json:"ptkId"`
	JenisSertifikasi string `json:"jenisSertifikasi"`
	NomorSertifikat  string `json:"nomorSertifikat"`
	TahunSertifikasi int    `json:"tahunSertifikasi"`
	DokumenURL       string `json:"dokumenUrl"`
}

type TeacherHistoryRepository interface {
	ListEducations(ctx context.Context, ptkID string) ([]EducationHistory, error)
	CreateEducation(ctx context.Context, edu *EducationHistory) error
	DeleteEducation(ctx context.Context, eduID string) error

	ListCertifications(ctx context.Context, ptkID string) ([]CertificationHistory, error)
	CreateCertification(ctx context.Context, cert *CertificationHistory) error
	DeleteCertification(ctx context.Context, certID string) error
}

type pgxTeacherHistoryRepository struct {
	db *pgxpool.Pool
}

func NewTeacherHistoryRepository(db *pgxpool.Pool) TeacherHistoryRepository {
	return &pgxTeacherHistoryRepository{db: db}
}

// === EDUCATION ===

func (r *pgxTeacherHistoryRepository) ListEducations(ctx context.Context, ptkID string) ([]EducationHistory, error) {
	query := `SELECT rwy_pend_formal_id, ptk_id, jenjang_pendidikan, COALESCE(jurusan,''), COALESCE(nama_sekolah,''), COALESCE(tahun_lulus,0), COALESCE(gelar,''), COALESCE(dokumen_url,'') 
	          FROM tabel_rwy_pend_formal WHERE ptk_id = $1 ORDER BY tahun_lulus DESC`
	rows, err := r.db.Query(ctx, query, ptkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []EducationHistory
	for rows.Next() {
		var e EducationHistory
		if err := rows.Scan(&e.ID, &e.PtkID, &e.Jenjang, &e.Jurusan, &e.NamaSekolah, &e.TahunLulus, &e.Gelar, &e.DokumenURL); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *pgxTeacherHistoryRepository) CreateEducation(ctx context.Context, edu *EducationHistory) error {
	if edu.ID == "" {
		edu.ID = uuid.New().String()
	}
	query := `INSERT INTO tabel_rwy_pend_formal (rwy_pend_formal_id, ptk_id, jenjang_pendidikan, jurusan, nama_sekolah, tahun_lulus, gelar, dokumen_url)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, edu.ID, edu.PtkID, edu.Jenjang, edu.Jurusan, edu.NamaSekolah, edu.TahunLulus, edu.Gelar, edu.DokumenURL)
	return err
}

func (r *pgxTeacherHistoryRepository) DeleteEducation(ctx context.Context, eduID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tabel_rwy_pend_formal WHERE rwy_pend_formal_id = $1`, eduID)
	return err
}

// === CERTIFICATION ===

func (r *pgxTeacherHistoryRepository) ListCertifications(ctx context.Context, ptkID string) ([]CertificationHistory, error) {
	query := `SELECT rwy_sertifikasi_id, ptk_id, jenis_sertifikasi, COALESCE(nomor_sertifikat,''), COALESCE(tahun_sertifikasi,0), COALESCE(dokumen_url,'') 
	          FROM tabel_rwy_sertifikasi WHERE ptk_id = $1 ORDER BY tahun_sertifikasi DESC`
	rows, err := r.db.Query(ctx, query, ptkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []CertificationHistory
	for rows.Next() {
		var c CertificationHistory
		if err := rows.Scan(&c.ID, &c.PtkID, &c.JenisSertifikasi, &c.NomorSertifikat, &c.TahunSertifikasi, &c.DokumenURL); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *pgxTeacherHistoryRepository) CreateCertification(ctx context.Context, cert *CertificationHistory) error {
	if cert.ID == "" {
		cert.ID = uuid.New().String()
	}
	query := `INSERT INTO tabel_rwy_sertifikasi (rwy_sertifikasi_id, ptk_id, jenis_sertifikasi, nomor_sertifikat, tahun_sertifikasi, dokumen_url)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, cert.ID, cert.PtkID, cert.JenisSertifikasi, cert.NomorSertifikat, cert.TahunSertifikasi, cert.DokumenURL)
	return err
}

func (r *pgxTeacherHistoryRepository) DeleteCertification(ctx context.Context, certID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tabel_rwy_sertifikasi WHERE rwy_sertifikasi_id = $1`, certID)
	return err
}
