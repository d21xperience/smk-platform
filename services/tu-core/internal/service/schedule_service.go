package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"sekolah-platform/services/tu-core/internal/repository"
)

type ScheduleRequest struct {
	ClassID    string `json:"class_id"`
	SubjectID  int    `json:"subject_id"`
	TeacherID  string `json:"teacher_id"`
	Hari       string `json:"hari"`
	JamMulai   string `json:"jam_mulai"`
	JamSelesai string `json:"jam_selesai"`
}

type ScheduleService interface {
	ListSchedules(ctx context.Context, schoolID, periodID string) ([]repository.Schedule, error)
	CreateSchedule(ctx context.Context, req ScheduleRequest, schoolID, periodID string) error
	DeleteSchedule(ctx context.Context, id string) error
}

type scheduleServiceImpl struct {
	scheduleRepo repository.ScheduleRepository
	db           *pgxpool.Pool
}

func NewScheduleService(sr repository.ScheduleRepository, db *pgxpool.Pool) ScheduleService {
	return &scheduleServiceImpl{scheduleRepo: sr, db: db}
}

func (s *scheduleServiceImpl) ListSchedules(ctx context.Context, schoolID, periodID string) ([]repository.Schedule, error) {
	return s.scheduleRepo.List(ctx, schoolID, periodID)
}

func (s *scheduleServiceImpl) CreateSchedule(ctx context.Context, req ScheduleRequest, schoolID, periodID string) error {
	var pembelajaranID string
	query := `
		SELECT pembelajaran_id::text 
		FROM tabel_pembelajaran 
		WHERE rombongan_belajar_id = $1 AND mata_pelajaran_id = $2 AND semester_id = $3
	`
	if req.TeacherID != "" {
		query += " AND ptk_terdaftar_id = $4"
		err := s.db.QueryRow(ctx, query, req.ClassID, req.SubjectID, periodID, req.TeacherID).Scan(&pembelajaranID)
		if err != nil {
			pembelajaranID = uuid.New().String()
			insertQuery := `INSERT INTO tabel_pembelajaran (pembelajaran_id, rombongan_belajar_id, mata_pelajaran_id, semester_id, ptk_terdaftar_id) VALUES ($1, $2, $3, $4, $5)`
			_, errInsert := s.db.Exec(ctx, insertQuery, pembelajaranID, req.ClassID, req.SubjectID, periodID, req.TeacherID)
			if errInsert != nil {
				return fmt.Errorf("gagal auto-create pembelajaran: %w", errInsert)
			}
		}
	} else {
		query += " AND ptk_terdaftar_id IS NULL"
		err := s.db.QueryRow(ctx, query, req.ClassID, req.SubjectID, periodID).Scan(&pembelajaranID)
		if err != nil {
			return fmt.Errorf("guru belum dipilih dan pembelajaran belum ada")
		}
	}

	schedule := &repository.Schedule{
		PembelajaranID: pembelajaranID,
		Hari:           req.Hari,
		JamMulai:       req.JamMulai,
		JamSelesai:     req.JamSelesai,
	}
	return s.scheduleRepo.Create(ctx, schedule, schoolID, periodID)
}

func (s *scheduleServiceImpl) DeleteSchedule(ctx context.Context, id string) error {
	return s.scheduleRepo.Delete(ctx, id)
}
