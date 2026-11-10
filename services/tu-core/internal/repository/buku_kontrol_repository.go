package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MasterGolongan struct {
	ID               int    `json:"id"`
	KodeGolongan     string `json:"kodeGolongan"`
	NamaPangkat      string `json:"namaPangkat"`
	IntervalKGBTahun int    `json:"intervalKgbTahun"`
}

type BukuKontrol struct {
	ID                  string     `json:"id"`
	PtkID               string     `json:"ptkId"`
	NamaPegawai         string     `json:"namaPegawai"`
	NIP                 string     `json:"nip"`
	GolonganSekarangID  int        `json:"golonganSekarangId"`
	KodeGolongan        string     `json:"kodeGolongan"`
	NamaPangkat         string     `json:"namaPangkat"`
	TmtPangkatSekarang  *time.Time `json:"tmtPangkatSekarang"`
	TmtGajiBerkalaSekarang *time.Time `json:"tmtGajiBerkalaSekarang"`
	TanggalKGBBerikut   *time.Time `json:"tanggalKgbBerikut"`
	TanggalPangkatBerikut *time.Time `json:"tanggalPangkatBerikut"`
	TanggalPensiun      *time.Time `json:"tanggalPensiun"`
	StatusMonitoring    string     `json:"statusMonitoring"`
	SelisihHari         int        `json:"selisihHari"`
}

type BukuKontrolRepository interface {
	ListGolongan(ctx context.Context) ([]MasterGolongan, error)
	ListBukuKontrol(ctx context.Context) ([]BukuKontrol, error)
	GetPegawaiForInit(ctx context.Context) ([]map[string]interface{}, error)
	CreateBukuKontrol(ctx context.Context, ptkID string, golID int, tmtPangkat, tmtGaji time.Time, tglKGB, tglPangkat, tglPensiun time.Time, status string) error
	GetBukuKontrolByPtkID(ctx context.Context, ptkID string) (*BukuKontrol, error)
}

type pgxBukuKontrolRepository struct {
	db *pgxpool.Pool
}

func NewBukuKontrolRepository(db *pgxpool.Pool) BukuKontrolRepository {
	return &pgxBukuKontrolRepository{db: db}
}

func (r *pgxBukuKontrolRepository) ListGolongan(ctx context.Context) ([]MasterGolongan, error) {
	query := `SELECT id, kode_golongan, nama_pangkat, interval_kgb_tahun FROM master_golongan ORDER BY id ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []MasterGolongan
	for rows.Next() {
		var g MasterGolongan
		if err := rows.Scan(&g.ID, &g.KodeGolongan, &g.NamaPangkat, &g.IntervalKGBTahun); err != nil {
			return nil, err
		}
		list = append(list, g)
	}
	return list, nil
}

func (r *pgxBukuKontrolRepository) ListBukuKontrol(ctx context.Context) ([]BukuKontrol, error) {
	query := `
		SELECT 
			bk.id::text, bk.ptk_id::text, 
			COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
			COALESCE(tp.nip, ''),
			COALESCE(bk.golongan_sekarang_id, 0), COALESCE(mg.kode_golongan, '-'), COALESCE(mg.nama_pangkat, '-'),
			bk.tmt_pangkat_sekarang, bk.tmt_gaji_berkala_sekarang,
			bk.tanggal_kgb_berikut, bk.tanggal_pangkat_berikut, bk.tanggal_pensiun, bk.status_monitoring,
			LEAST(
				COALESCE(bk.tanggal_kgb_berikut - CURRENT_DATE, 99999),
				COALESCE(bk.tanggal_pangkat_berikut - CURRENT_DATE, 99999)
			) AS selisih_hari
		FROM buku_kontrol bk
		JOIN tabel_ptk tp ON bk.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		LEFT JOIN master_golongan mg ON bk.golongan_sekarang_id = mg.id
		WHERE tp.soft_delete = 0
		ORDER BY selisih_hari ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []BukuKontrol
	for rows.Next() {
		var b BukuKontrol
		if err := rows.Scan(
			&b.ID, &b.PtkID, &b.NamaPegawai, &b.NIP,
			&b.GolonganSekarangID, &b.KodeGolongan, &b.NamaPangkat,
			&b.TmtPangkatSekarang, &b.TmtGajiBerkalaSekarang,
			&b.TanggalKGBBerikut, &b.TanggalPangkatBerikut, &b.TanggalPensiun, &b.StatusMonitoring,
			&b.SelisihHari,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *pgxBukuKontrolRepository) GetPegawaiForInit(ctx context.Context) ([]map[string]interface{}, error) {
	// Mengambil pegawai yang belum ada di tabel buku_kontrol
	query := `
		SELECT tp.ptk_id::text, tp.nama, COALESCE(tp.nip, ''), 
		       COALESCE(ext.golongan_id, 1), COALESCE(ext.tmt_pangkat, CURRENT_DATE), COALESCE(ext.tmt_gaji_berkala, CURRENT_DATE),
		       tp.tanggal_lahir
		FROM tabel_ptk tp
		LEFT JOIN ptk_buku_induk_ext ext ON tp.ptk_id = ext.ptk_id
		WHERE tp.soft_delete = 0 
		  AND NOT EXISTS (SELECT 1 FROM buku_kontrol bk WHERE bk.ptk_id = tp.ptk_id)
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, nama, nip string
		var golID int
		var tmtPangkat, tmtGaji, tglLahir time.Time
		if err := rows.Scan(&id, &nama, &nip, &golID, &tmtPangkat, &tmtGaji, &tglLahir); err != nil {
			return nil, err
		}
		list = append(list, map[string]interface{}{
			"ptkId": id, "nama": nama, "nip": nip,
			"golonganId": golID, "tmtPangkat": tmtPangkat, "tmtGaji": tmtGaji, "tanggalLahir": tglLahir,
		})
	}
	return list, nil
}

func (r *pgxBukuKontrolRepository) CreateBukuKontrol(ctx context.Context, ptkID string, golID int, tmtPangkat, tmtGaji time.Time, tglKGB, tglPangkat, tglPensiun time.Time, status string) error {
	query := `
		INSERT INTO buku_kontrol (
			ptk_id, golongan_sekarang_id, tmt_pangkat_sekarang, tmt_gaji_berkala_sekarang,
			tanggal_kgb_berikut, tanggal_pangkat_berikut, tanggal_pensiun, status_monitoring
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (ptk_id) DO UPDATE SET
			golongan_sekarang_id = EXCLUDED.golongan_sekarang_id,
			tmt_pangkat_sekarang = EXCLUDED.tmt_pangkat_sekarang,
			tmt_gaji_berkala_sekarang = EXCLUDED.tmt_gaji_berkala_sekarang,
			tanggal_kgb_berikut = EXCLUDED.tanggal_kgb_berikut,
			tanggal_pangkat_berikut = EXCLUDED.tanggal_pangkat_berikut,
			tanggal_pensiun = EXCLUDED.tanggal_pensiun,
			status_monitoring = EXCLUDED.status_monitoring,
			last_calculated_at = NOW()
	`
	_, err := r.db.Exec(ctx, query, ptkID, golID, tmtPangkat, tmtGaji, tglKGB, tglPangkat, tglPensiun, status)
	return err
}

func (r *pgxBukuKontrolRepository) GetBukuKontrolByPtkID(ctx context.Context, ptkID string) (*BukuKontrol, error) {
	query := `
		SELECT 
			bk.id::text, bk.ptk_id::text, tp.nama, COALESCE(tp.nip, ''),
			COALESCE(bk.golongan_sekarang_id, 0), COALESCE(mg.kode_golongan, '-'), COALESCE(mg.nama_pangkat, '-'),
			bk.tmt_pangkat_sekarang, bk.tmt_gaji_berkala_sekarang,
			bk.tanggal_kgb_berikut, bk.tanggal_pangkat_berikut, bk.tanggal_pensiun, bk.status_monitoring,
			0
		FROM buku_kontrol bk
		JOIN tabel_ptk tp ON bk.ptk_id = tp.ptk_id
		LEFT JOIN master_golongan mg ON bk.golongan_sekarang_id = mg.id
		WHERE bk.ptk_id = $1
	`
	var b BukuKontrol
	err := r.db.QueryRow(ctx, query, ptkID).Scan(
		&b.ID, &b.PtkID, &b.NamaPegawai, &b.NIP,
		&b.GolonganSekarangID, &b.KodeGolongan, &b.NamaPangkat,
		&b.TmtPangkatSekarang, &b.TmtGajiBerkalaSekarang,
		&b.TanggalKGBBerikut, &b.TanggalPangkatBerikut, &b.TanggalPensiun, &b.StatusMonitoring,
		&b.SelisihHari,
	)
	if err != nil { return nil, err }
	return &b, nil
}
