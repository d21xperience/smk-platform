package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type School struct {
	SchoolID          string   `json:"SchoolID"`
	Name              string   `json:"Name"`
	NPSN              string   `json:"NPSN"`
	NSS               string   `json:"NSS"`
	Address           string   `json:"Address"`
	PostalCode        string   `json:"PostalCode"`
	Phone             string   `json:"Phone"`
	Fax               string   `json:"Fax"`
	Kelurahan         string   `json:"Kelurahan"`
	Kecamatan         string   `json:"Kecamatan"`
	KabKota           string   `json:"KabKota"`
	Propinsi          string   `json:"Propinsi"`
	Website           string   `json:"Website"`
	Email             string   `json:"Email"`
	PrincipalName     string   `json:"PrincipalName"`
	PrincipalNIP      string   `json:"PrincipalNIP"`
	PrincipalNIY      string   `json:"PrincipalNIY"`
	StatusKepemilikan int      `json:"StatusKepemilikan"`
	KodeAktivasi      string   `json:"KodeAktivasi"`
	Jenjang           string   `json:"Jenjang"`
	BentukPendidikan  int      `json:"BentukPendidikan"`
	KodeSekolah       string   `json:"KodeSekolah"`
	LamaPendidikan    int      `json:"LamaPendidikan"`
	Jurusans          []string `json:"Jurusans"`
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
		       status_kepemilikan_id, COALESCE(kode_aktivasi, ''), COALESCE(jenjang, ''), COALESCE(bentuk_pendidikan_id, 0),
		       COALESCE(kode_sekolah, ''), COALESCE(lama_pendidikan, 3)
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
		&s.KodeSekolah, &s.LamaPendidikan,
	)
	if err != nil {
		return nil, err
	}

	// Ambil data Jurusan dari tabel Pivot
	rows, err := r.db.Query(ctx, `SELECT jurusan_id FROM tabel_sekolah_jurusan WHERE sekolah_id = $1`, schoolID)
	if err == nil {
		defer rows.Close()
		s.Jurusans = []string{}
		for rows.Next() {
			var jid string
			if err := rows.Scan(&jid); err == nil {
				s.Jurusans = append(s.Jurusans, jid)
			}
		}
	}

	return &s, nil
}

func (r *schoolRepository) Update(ctx context.Context, s *School) error {
	// PERBAIKAN: Memulai transaksi
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// PERBAIKAN: Rollback dengan context
	defer tx.Rollback(ctx)

	query := `
		UPDATE tabel_sekolah SET
			nama = $1, npsn = $2, nss = $3, alamat = $4, kd_pos = $5, telepon = $6, fax = $7,
			kelurahan = $8, kecamatan = $9, kab_kota = $10, propinsi = $11, website = $12, email = $13,
			nm_kepsek = $14, nip_kepsek = $15, niy_kepsek = $16, status_kepemilikan_id = $17,
			kode_aktivasi = $18, jenjang = $19, bentuk_pendidikan_id = $20, 
			kode_sekolah = $21, lama_pendidikan = $22
		WHERE sekolah_id = $23
	`
	_, err = tx.Exec(ctx, query,
		s.Name, s.NPSN, s.NSS, s.Address, s.PostalCode, s.Phone, s.Fax,
		s.Kelurahan, s.Kecamatan, s.KabKota, s.Propinsi, s.Website, s.Email,
		s.PrincipalName, s.PrincipalNIP, s.PrincipalNIY, s.StatusKepemilikan,
		s.KodeAktivasi, s.Jenjang, s.BentukPendidikan, s.KodeSekolah, s.LamaPendidikan, s.SchoolID,
	)
	if err != nil {
		return err
	}

	// Hapus relasi jurusan lama
	_, err = tx.Exec(ctx, `DELETE FROM tabel_sekolah_jurusan WHERE sekolah_id = $1`, s.SchoolID)
	if err != nil {
		return err
	}

	// Masukkan relasi jurusan baru (bisa lebih dari satu)
	for _, jid := range s.Jurusans {
		_, err = tx.Exec(ctx, `INSERT INTO tabel_sekolah_jurusan (sekolah_id, jurusan_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, s.SchoolID, jid)
		if err != nil {
			return err
		}
	}

	// PERBAIKAN: Commit dengan context
	return tx.Commit(ctx)
}
