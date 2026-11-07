package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeaveRequest struct {
	ID              string    `json:"id"`
	PtkID           string    `json:"ptkId"`
	TeacherName     string    `json:"teacherName"`
	JenisCuti       string    `json:"jenisCuti"`
	TanggalMulai    time.Time `json:"tanggalMulai"`
	TanggalSelesai  time.Time `json:"tanggalSelesai"`
	Alasan          string    `json:"alasan"`
	Status          string    `json:"status"`
	CatatanApproval string    `json:"catatanApproval"`
	CreatedAt       time.Time `json:"createdAt"`
}

type PtkLeaveRepository interface {
	ListLeaves(ctx context.Context) ([]LeaveRequest, error)
	CreateLeave(ctx context.Context, req *LeaveRequest) error
	UpdateStatus(ctx context.Context, id string, status string, note string) error
	DeleteLeave(ctx context.Context, id string) error
}

type pgxPtkLeaveRepository struct {
	db *pgxpool.Pool
}

func NewPtkLeaveRepository(db *pgxpool.Pool) PtkLeaveRepository {
	return &pgxPtkLeaveRepository{db: db}
}

func (r *pgxPtkLeaveRepository) ListLeaves(ctx context.Context) ([]LeaveRequest, error) {
	query := `
		SELECT l.id::text, l.ptk_id::text, 
		       COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
		       l.jenis_cuti, l.tanggal_mulai, l.tanggal_selesai, COALESCE(l.alasan, ''), l.status, COALESCE(l.catatan_approval, ''), l.created_at
		FROM ptk_leaves l
		JOIN tabel_ptk tp ON l.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		ORDER BY l.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []LeaveRequest
	for rows.Next() {
		var l LeaveRequest
		if err := rows.Scan(&l.ID, &l.PtkID, &l.TeacherName, &l.JenisCuti, &l.TanggalMulai, &l.TanggalSelesai, &l.Alasan, &l.Status, &l.CatatanApproval, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (r *pgxPtkLeaveRepository) CreateLeave(ctx context.Context, req *LeaveRequest) error {
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	query := `
		INSERT INTO ptk_leaves (id, ptk_id, jenis_cuti, tanggal_mulai, tanggal_selesai, alasan, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending')
	`
	_, err := r.db.Exec(ctx, query, req.ID, req.PtkID, req.JenisCuti, req.TanggalMulai, req.TanggalSelesai, req.Alasan)
	return err
}

func (r *pgxPtkLeaveRepository) UpdateStatus(ctx context.Context, id string, status string, note string) error {
	query := `
		UPDATE ptk_leaves 
		SET status = $1, catatan_approval = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, status, note, id)
	return err
}

func (r *pgxPtkLeaveRepository) DeleteLeave(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM ptk_leaves WHERE id = $1`, id)
	return err
}
