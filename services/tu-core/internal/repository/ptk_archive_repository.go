package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Archive struct {
	ID               string    `json:"id"`
	PtkID            string    `json:"ptkId"`
	TeacherName      string    `json:"teacherName"`
	JenisDokumen     string    `json:"jenisDokumen"`
	NomorDokumen     string    `json:"nomorDokumen"`
	TanggalTerbit    time.Time `json:"tanggalTerbit"`
	TanggalKadaluarsa *time.Time `json:"tanggalKadaluarsa"`
	Catatan          string    `json:"catatan"`
	FileURL          string    `json:"fileUrl"`
	CreatedAt        time.Time `json:"createdAt"`
}

type PtkArchiveRepository interface {
	ListArchives(ctx context.Context) ([]Archive, error)
	CreateArchive(ctx context.Context, req *Archive) error
	UpdateArchive(ctx context.Context, req *Archive) error
	DeleteArchive(ctx context.Context, id string) error
}

type pgxPtkArchiveRepository struct {
	db *pgxpool.Pool
}

func NewPtkArchiveRepository(db *pgxpool.Pool) PtkArchiveRepository {
	return &pgxPtkArchiveRepository{db: db}
}

func (r *pgxPtkArchiveRepository) ListArchives(ctx context.Context) ([]Archive, error) {
	query := `
		SELECT a.id::text, a.ptk_id::text, 
		       COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
		       a.jenis_dokumen, COALESCE(a.nomor_dokumen, ''), 
		       COALESCE(a.tanggal_terbit, '1970-01-01'::date), 
		       a.tanggal_kadaluarsa, 
		       COALESCE(a.catatan, ''), a.file_url, a.created_at
		FROM ptk_archives a
		JOIN tabel_ptk tp ON a.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		ORDER BY a.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Archive
	for rows.Next() {
		var a Archive
		if err := rows.Scan(&a.ID, &a.PtkID, &a.TeacherName, &a.JenisDokumen, &a.NomorDokumen, &a.TanggalTerbit, &a.TanggalKadaluarsa, &a.Catatan, &a.FileURL, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *pgxPtkArchiveRepository) CreateArchive(ctx context.Context, req *Archive) error {
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	query := `
		INSERT INTO ptk_archives (id, ptk_id, jenis_dokumen, nomor_dokumen, tanggal_terbit, tanggal_kadaluarsa, catatan, file_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	var tglTerbit any = nil
	if !req.TanggalTerbit.IsZero() { tglTerbit = req.TanggalTerbit }

	_, err := r.db.Exec(ctx, query, req.ID, req.PtkID, req.JenisDokumen, req.NomorDokumen, tglTerbit, req.TanggalKadaluarsa, req.Catatan, req.FileURL)
	return err
}

func (r *pgxPtkArchiveRepository) UpdateArchive(ctx context.Context, req *Archive) error {
	query := `
		UPDATE ptk_archives SET 
			jenis_dokumen = $1, nomor_dokumen = $2, tanggal_terbit = $3, 
			tanggal_kadaluarsa = $4, catatan = $5, file_url = $6
		WHERE id = $7
	`
	var tglTerbit any = nil
	if !req.TanggalTerbit.IsZero() { tglTerbit = req.TanggalTerbit }

	_, err := r.db.Exec(ctx, query, req.JenisDokumen, req.NomorDokumen, tglTerbit, req.TanggalKadaluarsa, req.Catatan, req.FileURL, req.ID)
	return err
}

func (r *pgxPtkArchiveRepository) DeleteArchive(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM ptk_archives WHERE id = $1`, id)
	return err
}
