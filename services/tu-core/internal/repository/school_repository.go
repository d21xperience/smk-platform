package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type School struct {
	SchoolID          string `json:"SchoolID"`
	Name              string `json:"Name"`
	NPSN              string `json:"NPSN"`
	NSS               string `json:"NSS"`
	Address           string `json:"Address"`
	PostalCode        string `json:"PostalCode"`
	Phone             string `json:"Phone"`
	Fax               string `json:"Fax"`
	Kelurahan         string `json:"Kelurahan"`
	Kecamatan         string `json:"Kecamatan"`
	KabKota           string `json:"KabKota"`
	Propinsi          string `json:"Propinsi"`
	Website           string `json:"Website"`
	Email             string `json:"Email"`
	PrincipalName     string `json:"PrincipalName"`
	PrincipalNIP      string `json:"PrincipalNIP"`
	PrincipalNIY      string `json:"PrincipalNIY"`
	StatusKepemilikan int    `json:"StatusKepemilikan"`
	KodeAktivasi      string `json:"KodeAktivasi"`
	Jenjang           string `json:"Jenjang"`
	BentukPendidikan  int    `json:"BentukPendidikan"`
}

type SchoolRepository interface {
	GetByID(ctx context.Context, schoolID string) (*School, error)
	Update(ctx context.Context, school *School) error
}

type schoolRepository struct {
	db *pgxpool.Pool
}

func NewSchoolRepository(db *pgxpool.Pool) SchoolRepository {
	return &schoolRepository{db: db}
}

func (r *schoolRepository) GetByID(ctx context.Context, schoolID string) (*School, error) {
	query := `
		SELECT sekolah_id::text, nama, COALESCE(npsn, ''), COALESCE(nss, ''), COALESCE(alamat, ''),
		       COALESCE(kd_pos, ''), COALESCE(telepon, ''), COALESCE(fax, ''), COALESCE(kelurahan, ''),
		       COALESCE(kecamatan, ''), COALESCE(kab_kota, ''), COALESCE(propinsi, ''), COALESCE(website, ''),
		       COALESCE(email, ''), COALESCE(nm_kepsek, ''), COALESCE(nip_kepsek, ''), COALESCE(niy_kepsek, ''),
		       status_kepemilikan_id, COALESCE(kode_aktivasi, ''), COALESCE(jenjang, ''), COALESCE(bentuk_pendidikan_id, 0)
		FROM tabel_sekolah
		WHERE sekolah_id = $1
	`
	var s School
	err := r.db.QueryRow(ctx, query, schoolID).Scan(
		&s.SchoolID, &s.Name, &s.NPSN, &s.NSS, &s.Address,
		&s.PostalCode, &s.Phone, &s.Fax, &s.Kelurahan,
		&s.Kecamatan, &s.KabKota, &s.Propinsi, &s.Website,
		&s.Email, &s.PrincipalName, &s.PrincipalNIP, &s.PrincipalNIY,
		&s.StatusKepemilikan, &s.KodeAktivasi, &s.Jenjang, &s.BentukPendidikan,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *schoolRepository) Update(ctx context.Context, s *School) error {
	query := `
		UPDATE tabel_sekolah SET
			nama = $1, npsn = $2, nss = $3, alamat = $4, kd_pos = $5, telepon = $6, fax = $7,
			kelurahan = $8, kecamatan = $9, kab_kota = $10, propinsi = $11, website = $12, email = $13,
			nm_kepsek = $14, nip_kepsek = $15, niy_kepsek = $16, status_kepemilikan_id = $17,
			kode_aktivasi = $18, jenjang = $19, bentuk_pendidikan_id = $20
		WHERE sekolah_id = $21
	`
	_, err := r.db.Exec(ctx, query,
		s.Name, s.NPSN, s.NSS, s.Address, s.PostalCode, s.Phone, s.Fax,
		s.Kelurahan, s.Kecamatan, s.KabKota, s.Propinsi, s.Website, s.Email,
		s.PrincipalName, s.PrincipalNIP, s.PrincipalNIY, s.StatusKepemilikan,
		s.KodeAktivasi, s.Jenjang, s.BentukPendidikan, s.SchoolID,
	)
	return err
}
