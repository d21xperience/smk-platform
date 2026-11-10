package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BukuIndukDetail struct {
	PtkID             string     `json:"ptkId"`
	Name              string     `json:"name"`
	NIP               string     `json:"nip"`
	NUPTK             string     `json:"nuptk"`
	NIK               string     `json:"nik"`
	Gender            string     `json:"gender"`
	BirthPlace        string     `json:"birthPlace"`
	BirthDate         *time.Time `json:"birthDate"`
	Address           string     `json:"address"`
	NoHP              string     `json:"noHp"`
	Email             string     `json:"email"`
	FotoURL           string     `json:"fotoUrl"`

	// Dari ptk_buku_induk_ext
	JenisPegawai      string     `json:"jenisPegawai"`
	StatusPegawai     string     `json:"statusPegawai"`
	TMTMasuk          *time.Time `json:"tmtMasuk"`
	UnitKerja         string     `json:"unitKerja"`
	Jabatan           string     `json:"jabatan"`
	PangkatGolongan   string     `json:"pangkatGolongan"`
	
	// Kelengkapan
	SkorKelengkapan   int        `json:"skorKelengkapan"`
	StatusKelengkapan string     `json:"statusKelengkapan"`
}

type BukuIndukRepository interface {
	ListBukuInduk(ctx context.Context) ([]BukuIndukDetail, error)
	GetBukuIndukDetail(ctx context.Context, ptkID string) (*BukuIndukDetail, error)
}

type pgxBukuIndukRepository struct {
	db *pgxpool.Pool
}

func NewBukuIndukRepository(db *pgxpool.Pool) BukuIndukRepository {
	return &pgxBukuIndukRepository{db: db}
}

func (r *pgxBukuIndukRepository) ListBukuInduk(ctx context.Context) ([]BukuIndukDetail, error) {
	// Query dasar. Skor kelengkapan dihitung di level Go (Service) agar fleksibel
	query := `
		SELECT 
			tp.ptk_id::text, 
			COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
			COALESCE(tp.nip, ''), COALESCE(tp.nuptk, ''), COALESCE(tp.nik, ''), tp.jenis_kelamin,
			tp.tempat_lahir, tp.tanggal_lahir, COALESCE(tp.alamat_jalan, ''), COALESCE(tp.no_hp, ''), COALESCE(tp.email, ''), COALESCE(tp.foto_url, ''),
			COALESCE(ext.jenis_pegawai, 'Guru'), COALESCE(ext.status_pegawai, 'Honorer'), ext.tmt_masuk,
			COALESCE(ext.unit_kerja, ''), COALESCE(ext.jabatan, ''), COALESCE(ext.pangkat_golongan, '')
		FROM tabel_ptk tp
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		LEFT JOIN ptk_buku_induk_ext ext ON tp.ptk_id = ext.ptk_id
		WHERE tp.soft_delete = 0
		ORDER BY nama_lengkap ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []BukuIndukDetail
	for rows.Next() {
		var d BukuIndukDetail
		if err := rows.Scan(
			&d.PtkID, &d.Name, &d.NIP, &d.NUPTK, &d.NIK, &d.Gender,
			&d.BirthPlace, &d.BirthDate, &d.Address, &d.NoHP, &d.Email, &d.FotoURL,
			&d.JenisPegawai, &d.StatusPegawai, &d.TMTMasuk, &d.UnitKerja, &d.Jabatan, &d.PangkatGolongan,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *pgxBukuIndukRepository) GetBukuIndukDetail(ctx context.Context, ptkID string) (*BukuIndukDetail, error) {
	query := `
		SELECT 
			tp.ptk_id::text, 
			COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
			COALESCE(tp.nip, ''), COALESCE(tp.nuptk, ''), COALESCE(tp.nik, ''), tp.jenis_kelamin,
			tp.tempat_lahir, tp.tanggal_lahir, COALESCE(tp.alamat_jalan, ''), COALESCE(tp.no_hp, ''), COALESCE(tp.email, ''), COALESCE(tp.foto_url, ''),
			COALESCE(ext.jenis_pegawai, 'Guru'), COALESCE(ext.status_pegawai, 'Honorer'), ext.tmt_masuk,
			COALESCE(ext.unit_kerja, ''), COALESCE(ext.jabatan, ''), COALESCE(ext.pangkat_golongan, '')
		FROM tabel_ptk tp
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		LEFT JOIN ptk_buku_induk_ext ext ON tp.ptk_id = ext.ptk_id
		WHERE tp.ptk_id = $1 AND tp.soft_delete = 0 LIMIT 1
	`
	var d BukuIndukDetail
	err := r.db.QueryRow(ctx, query, ptkID).Scan(
		&d.PtkID, &d.Name, &d.NIP, &d.NUPTK, &d.NIK, &d.Gender,
		&d.BirthPlace, &d.BirthDate, &d.Address, &d.NoHP, &d.Email, &d.FotoURL,
		&d.JenisPegawai, &d.StatusPegawai, &d.TMTMasuk, &d.UnitKerja, &d.Jabatan, &d.PangkatGolongan,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
