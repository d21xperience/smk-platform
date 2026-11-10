package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsulanKenaikan struct {
	ID                string     `json:"id"`
	NomorUsulan       string     `json:"nomorUsulan"`
	PtkID             string     `json:"ptkId"`
	NamaPegawai       string     `json:"namaPegawai"`
	NIP               string     `json:"nip"`
	Jenis             string     `json:"jenis"`
	Status            string     `json:"status"`
	GolonganSekarang  string     `json:"golonganSekarang"`
	GolonganBaru      string     `json:"golonganBaru"`
	TmtBaru           *time.Time `json:"tmtBaru"`
	NomorSKTerbit     string     `json:"nomorSkTerbit"`
	CreatedAt         time.Time  `json:"createdAt"`
}

type SaranUsulan struct {
	PtkID           string `json:"ptkId"`
	NamaPegawai     string `json:"namaPegawai"`
	NIP             string `json:"nip"`
	JenisUsulan     string `json:"jenisUsulan"`
	StatusMonitoring string `json:"statusMonitoring"`
	SelisihHari     int    `json:"selisihHari"`
	GolonganSekarang string `json:"golonganSekarang"`
}

type UsulanKenaikanRepository interface {
	ListUsulan(ctx context.Context) ([]UsulanKenaikan, error)
	GetSaranUsulan(ctx context.Context) ([]SaranUsulan, error)
	CreateUsulan(ctx context.Context, ptkID, jenis, nomorUsulan, userID string, golSekarangID, golBaruID int, tmtBaru time.Time, catatan string) (string, error)
	UpdateStatus(ctx context.Context, id, status, catatan string) error
}

type pgxUsulanKenaikanRepository struct {
	db *pgxpool.Pool
}

func NewUsulanKenaikanRepository(db *pgxpool.Pool) UsulanKenaikanRepository {
	return &pgxUsulanKenaikanRepository{db: db}
}

func (r *pgxUsulanKenaikanRepository) ListUsulan(ctx context.Context) ([]UsulanKenaikan, error) {
	query := `
		SELECT 
			uk.id::text, uk.nomor_usulan, uk.ptk_id::text, 
			COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
			COALESCE(tp.nip, ''), uk.jenis, uk.status,
			COALESCE(mg1.kode_golongan, '-') AS gol_sekarang,
			COALESCE(mg2.kode_golongan, '-') AS gol_baru,
			uk.tmt_baru, COALESCE(uk.nomor_sk_terbit, ''), uk.created_at
		FROM usulan_kenaikan uk
		JOIN tabel_ptk tp ON uk.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		LEFT JOIN master_golongan mg1 ON uk.golongan_sekarang_id = mg1.id
		LEFT JOIN master_golongan mg2 ON uk.golongan_baru_id = mg2.id
		ORDER BY uk.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]UsulanKenaikan, 0)
	for rows.Next() {
		var u UsulanKenaikan
		if err := rows.Scan(
			&u.ID, &u.NomorUsulan, &u.PtkID, &u.NamaPegawai, &u.NIP, &u.Jenis, &u.Status,
			&u.GolonganSekarang, &u.GolonganBaru, &u.TmtBaru, &u.NomorSKTerbit, &u.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

func (r *pgxUsulanKenaikanRepository) GetSaranUsulan(ctx context.Context) ([]SaranUsulan, error) {
	query := `
		SELECT 
			bk.ptk_id::text, 
			COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
			COALESCE(tp.nip, ''),
			CASE 
				WHEN bk.tanggal_pangkat_berikut <= bk.tanggal_kgb_berikut THEN 'KENAIKAN_PANGKAT'
				ELSE 'KGB'
			END AS jenis_usulan,
			bk.status_monitoring,
			LEAST(
				COALESCE(bk.tanggal_kgb_berikut - CURRENT_DATE, 99999),
				COALESCE(bk.tanggal_pangkat_berikut - CURRENT_DATE, 99999)
			) AS selisih_hari,
			COALESCE(mg.kode_golongan, '-')
		FROM buku_kontrol bk
		JOIN tabel_ptk tp ON bk.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		LEFT JOIN master_golongan mg ON bk.golongan_sekarang_id = mg.id
		WHERE bk.status_monitoring IN ('URGENT', 'SEGERA', 'PERLU_PERSIAPAN', 'OVERDUE')
		  AND NOT EXISTS (
			  SELECT 1 FROM usulan_kenaikan uk 
			  WHERE uk.ptk_id = bk.ptk_id 
			  AND uk.status NOT IN ('DITOLAK', 'SK_TERBIT')
		  )
		ORDER BY selisih_hari ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]SaranUsulan, 0)
	for rows.Next() {
		var s SaranUsulan
		if err := rows.Scan(
			&s.PtkID, &s.NamaPegawai, &s.NIP, &s.JenisUsulan, &s.StatusMonitoring, &s.SelisihHari, &s.GolonganSekarang,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *pgxUsulanKenaikanRepository) CreateUsulan(ctx context.Context, ptkID, jenis, nomorUsulan, userID string, golSekarangID, golBaruID int, tmtBaru time.Time, catatan string) (string, error) {
	var id string
	query := `
		INSERT INTO usulan_kenaikan (
			nomor_usulan, ptk_id, jenis, status, 
			golongan_sekarang_id, golongan_baru_id, tmt_baru, 
			catatan_ktu, created_by
		) VALUES ($1, $2, $3, 'DRAFT', $4, $5, $6, $7, $8)
		RETURNING id::text
	`
	err := r.db.QueryRow(ctx, query, nomorUsulan, ptkID, jenis, golSekarangID, golBaruID, tmtBaru, catatan, userID).Scan(&id)
	return id, err
}

func (r *pgxUsulanKenaikanRepository) UpdateStatus(ctx context.Context, id, status, catatan string) error {
	query := `UPDATE usulan_kenaikan SET status = $1, catatan_kepsek = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, status, catatan, id)
	return err
}
