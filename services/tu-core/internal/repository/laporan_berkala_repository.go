package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LaporanBerkala struct {
	ID              string          `json:"id"`
	SchoolID        string          `json:"schoolId"`
	Tahun           int             `json:"tahun"`
	Triwulan        int             `json:"triwulan"`
	TanggalMulai    time.Time       `json:"tanggalMulai"`
	TanggalSelesai  time.Time       `json:"tanggalSelesai"`
	Status          string          `json:"status"`
	DataJSON        json.RawMessage `json:"dataJson"`
	PDFUrl          string          `json:"pdfUrl"`
	ExcelDinasURL   string          `json:"excelDinasUrl"`
	ExcelYayasanURL string          `json:"excelYayasanUrl"`
	RingkasanURL    string          `json:"ringkasanUrl"`
	SuratURL        string          `json:"suratUrl"`
	CreatedAt       time.Time       `json:"createdAt"`
}

type LaporanBerkalaRepository interface {
	CreateDraft(ctx context.Context, l *LaporanBerkala) error
	GetByID(ctx context.Context, id string) (*LaporanBerkala, error)
	List(ctx context.Context, schoolID string) ([]LaporanBerkala, error)
	UpdateStatus(ctx context.Context, id, status, userID, notes string) error
	UpdateDocuments(ctx context.Context, id string, pdfURL, excelDinas, excelYayasan, ringkasan, surat string) error
	LogTimeline(ctx context.Context, laporanID, statusDari, statusKe, userID, catatan string) error
}

type pgxLaporanBerkalaRepository struct {
	db *pgxpool.Pool
}

func NewLaporanBerkalaRepository(db *pgxpool.Pool) LaporanBerkalaRepository {
	return &pgxLaporanBerkalaRepository{db: db}
}

func (r *pgxLaporanBerkalaRepository) CreateDraft(ctx context.Context, l *LaporanBerkala) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	query := `
		INSERT INTO laporan_berkala (
			id, school_id, tahun, triwulan, tanggal_mulai, tanggal_selesai, status, data_json
		) VALUES ($1, $2, $3, $4, $5, $6, 'DRAFT', $7)
	`
	_, err := r.db.Exec(ctx, query, l.ID, l.SchoolID, l.Tahun, l.Triwulan, l.TanggalMulai, l.TanggalSelesai, l.DataJSON)
	return err
}

func (r *pgxLaporanBerkalaRepository) GetByID(ctx context.Context, id string) (*LaporanBerkala, error) {
	query := `SELECT id, school_id, tahun, triwulan, tanggal_mulai, tanggal_selesai, status, COALESCE(data_json, '{}'), COALESCE(pdf_url, ''), COALESCE(excel_dinas_url, ''), COALESCE(excel_yayasan_url, ''), COALESCE(ringkasan_url, ''), COALESCE(surat_url, ''), created_at FROM laporan_berkala WHERE id = $1`
	var l LaporanBerkala
	err := r.db.QueryRow(ctx, query, id).Scan(
		&l.ID, &l.SchoolID, &l.Tahun, &l.Triwulan, &l.TanggalMulai, &l.TanggalSelesai,
		&l.Status, &l.DataJSON, &l.PDFUrl, &l.ExcelDinasURL, &l.ExcelYayasanURL, &l.RingkasanURL, &l.SuratURL, &l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *pgxLaporanBerkalaRepository) List(ctx context.Context, schoolID string) ([]LaporanBerkala, error) {
	query := `SELECT id, school_id, tahun, triwulan, tanggal_mulai, tanggal_selesai, status, created_at FROM laporan_berkala WHERE school_id = $1 ORDER BY tahun DESC, triwulan DESC`
	rows, err := r.db.Query(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []LaporanBerkala
	for rows.Next() {
		var l LaporanBerkala
		if err := rows.Scan(&l.ID, &l.SchoolID, &l.Tahun, &l.Triwulan, &l.TanggalMulai, &l.TanggalSelesai, &l.Status, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (r *pgxLaporanBerkalaRepository) UpdateStatus(ctx context.Context, id, status, userID, notes string) error {
	query := `UPDATE laporan_berkala SET status = $1, updated_at = NOW() WHERE id = $2`
	switch status {
	case "REVIEWED":
		query = `UPDATE laporan_berkala SET status = $1, catatan_ktu = $3, reviewed_by = $4, updated_at = NOW() WHERE id = $2`
		_, err := r.db.Exec(ctx, query, status, id, notes, userID)
		return err
	case "APPROVED":
		query = `UPDATE laporan_berkala SET status = $1, catatan_approve = $3, approved_by = $4, updated_at = NOW() WHERE id = $2`
		_, err := r.db.Exec(ctx, query, status, id, notes, userID)
		return err
	case "SENT":
		query = `UPDATE laporan_berkala SET status = $1, sent_by = $4, sent_at = NOW(), updated_at = NOW() WHERE id = $2`
		_, err := r.db.Exec(ctx, query, status, id, notes, userID)
		return err
	}
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *pgxLaporanBerkalaRepository) UpdateDocuments(ctx context.Context, id string, pdfURL, excelDinas, excelYayasan, ringkasan, surat string) error {
	query := `
		UPDATE laporan_berkala SET
			pdf_url = COALESCE(NULLIF($1, ''), pdf_url),
			excel_dinas_url = COALESCE(NULLIF($2, ''), excel_dinas_url),
			excel_yayasan_url = COALESCE(NULLIF($3, ''), excel_yayasan_url),
			ringkasan_url = COALESCE(NULLIF($4, ''), ringkasan_url),
			surat_url = COALESCE(NULLIF($5, ''), surat_url),
			updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, pdfURL, excelDinas, excelYayasan, ringkasan, surat, id)
	return err
}

func (r *pgxLaporanBerkalaRepository) LogTimeline(ctx context.Context, laporanID, statusDari, statusKe, userID, catatan string) error {
	query := `INSERT INTO laporan_timeline (laporan_id, status_dari, status_ke, action_by, catatan) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, laporanID, statusDari, statusKe, userID, catatan)
	return err
}
