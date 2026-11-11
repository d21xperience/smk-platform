package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MasterLayanan struct {
	Kode          string `json:"kode"`
	Nama          string `json:"nama"`
	Kategori      string `json:"kategori"`
	Deskripsi     string `json:"deskripsi"`
	SLABulan      int    `json:"slaHari"`
	TemplateSurat string `json:"templateSurat"`
}

type PengajuanLayanan struct {
	ID              string     `json:"id"`
	NomorPengajuan  string     `json:"nomorPengajuan"`
	PtkID           string     `json:"ptkId"`
	NamaPegawai     string     `json:"namaPegawai"`
	NIP             string     `json:"nip"`
	LayananKode     string     `json:"layananKode"`
	LayananNama     string     `json:"layananNama"`
	Status          string     `json:"status"`
	TujuanSurat     string     `json:"tujuanSurat"`
	Kepada          string     `json:"kepada"`
	Keterangan      string     `json:"keterangan"`
	EstimasiSelesai *time.Time `json:"estimasiSelesai"`
	RequestedAt     time.Time  `json:"requestedAt"`
}

type PelayananRepository interface {
	ListKatalog(ctx context.Context) ([]MasterLayanan, error)
	GetKatalogByKode(ctx context.Context, kode string) (*MasterLayanan, error)
	CreatePengajuan(ctx context.Context, ptkID, kode, tujuan, kepada, ket, dokURL, nomor string, estimasi time.Time) (string, error)
	ListPengajuan(ctx context.Context, status string) ([]PengajuanLayanan, error)
	GetPengajuanByID(ctx context.Context, id string) (*PengajuanLayanan, error)
	UpdateStatus(ctx context.Context, id, status, catatan, userID string) error
	AddTimeline(ctx context.Context, pengajuanID, dari, ke, userID, catatan string) error
}

type pgxPelayananRepository struct {
	db *pgxpool.Pool
}

func NewPelayananRepository(db *pgxpool.Pool) PelayananRepository {
	return &pgxPelayananRepository{db: db}
}

func (r *pgxPelayananRepository) ListKatalog(ctx context.Context) ([]MasterLayanan, error) {
	query := `SELECT kode, nama, kategori, COALESCE(deskripsi, ''), sla_hari, template_surat FROM master_layanan WHERE is_active = true ORDER BY kategori, nama`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]MasterLayanan, 0)
	for rows.Next() {
		var m MasterLayanan
		if err := rows.Scan(&m.Kode, &m.Nama, &m.Kategori, &m.Deskripsi, &m.SLABulan, &m.TemplateSurat); err == nil {
			list = append(list, m)
		}
	}
	return list, nil
}

func (r *pgxPelayananRepository) GetKatalogByKode(ctx context.Context, kode string) (*MasterLayanan, error) {
	query := `SELECT kode, nama, kategori, deskripsi, sla_hari, template_surat FROM master_layanan WHERE kode = $1`
	var m MasterLayanan
	err := r.db.QueryRow(ctx, query, kode).Scan(&m.Kode, &m.Nama, &m.Kategori, &m.Deskripsi, &m.SLABulan, &m.TemplateSurat)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *pgxPelayananRepository) CreatePengajuan(ctx context.Context, ptkID, kode, tujuan, kepada, ket, dokURL, nomor string, estimasi time.Time) (string, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO pelayanan_pengajuan (
			id, nomor_pengajuan, ptk_id, layanan_kode, status,
			tujuan_surat, kepada, keterangan, dokumen_pendukung_url, estimasi_selesai
		) VALUES ($1, $2, $3, $4, 'SUBMITTED', $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query, id, nomor, ptkID, kode, tujuan, kepada, ket, dokURL, estimasi)
	return id, err
}

func (r *pgxPelayananRepository) ListPengajuan(ctx context.Context, status string) ([]PengajuanLayanan, error) {
	query := `
		SELECT p.id::text, p.nomor_pengajuan, p.ptk_id::text,
		       COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, ''),
		       COALESCE(tp.nip, ''), p.layanan_kode, ml.nama, p.status,
		       COALESCE(p.tujuan_surat, ''), COALESCE(p.kepada, ''), COALESCE(p.keterangan, ''),
		       p.estimasi_selesai, p.requested_at
		FROM pelayanan_pengajuan p
		JOIN tabel_ptk tp ON p.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		JOIN master_layanan ml ON p.layanan_kode = ml.kode
		WHERE ($1 = '' OR p.status = $1)
		ORDER BY p.requested_at DESC
	`
	rows, err := r.db.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]PengajuanLayanan, 0)
	for rows.Next() {
		var p PengajuanLayanan
		if err := rows.Scan(
			&p.ID, &p.NomorPengajuan, &p.PtkID, &p.NamaPegawai, &p.NIP,
			&p.LayananKode, &p.LayananNama, &p.Status,
			&p.TujuanSurat, &p.Kepada, &p.Keterangan, &p.EstimasiSelesai, &p.RequestedAt,
		); err == nil {
			list = append(list, p)
		}
	}
	return list, nil
}

func (r *pgxPelayananRepository) GetPengajuanByID(ctx context.Context, id string) (*PengajuanLayanan, error) {
	query := `
		SELECT p.id::text, p.nomor_pengajuan, p.ptk_id::text, tp.nama, COALESCE(tp.nip, ''),
		       p.layanan_kode, ml.nama, p.status,
		       COALESCE(p.tujuan_surat, ''), COALESCE(p.kepada, ''), COALESCE(p.keterangan, ''),
		       p.estimasi_selesai, p.requested_at
		FROM pelayanan_pengajuan p
		JOIN tabel_ptk tp ON p.ptk_id = tp.ptk_id
		JOIN master_layanan ml ON p.layanan_kode = ml.kode
		WHERE p.id = $1
	`
	var p PengajuanLayanan
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.NomorPengajuan, &p.PtkID, &p.NamaPegawai, &p.NIP,
		&p.LayananKode, &p.LayananNama, &p.Status,
		&p.TujuanSurat, &p.Kepada, &p.Keterangan, &p.EstimasiSelesai, &p.RequestedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pgxPelayananRepository) UpdateStatus(ctx context.Context, id, status, catatan, userID string) error {
	var query string
	switch status {
	case "REJECTED":
		query = `UPDATE pelayanan_pengajuan SET status = $1, alasan_penolakan = $2 WHERE id = $3`
	case "REVISION":
		query = `UPDATE pelayanan_pengajuan SET status = $1, catatan_verifikasi = $2 WHERE id = $3`
	default:
		query = `UPDATE pelayanan_pengajuan SET status = $1 WHERE id = $2`
		_, err := r.db.Exec(ctx, query, status, id)
		return err
	}
	_, err := r.db.Exec(ctx, query, status, catatan, id)
	return err
}

func (r *pgxPelayananRepository) AddTimeline(ctx context.Context, pengajuanID, dari, ke, userID, catatan string) error {
	query := `
		INSERT INTO pelayanan_timeline (pengajuan_id, status_dari, status_ke, action_by, catatan)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, pengajuanID, dari, ke, userID, catatan)
	return err
}
