package repository

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/agenda/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AgendaFilter struct {
	SchoolID string
	Jenis    string // IN, OUT, BOTH
	Bulan    int
	Tahun    int
}

type AgendaRepository interface {
	GenerateNomorAgenda(ctx context.Context, schoolID, jenis string, tahun int) (string, int, error)
	CreateAgendaEntry(ctx context.Context, entry *models.AgendaEntry) error
	GetByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) (*models.AgendaEntry, error)
	ListAgenda(ctx context.Context, filter AgendaFilter) ([]map[string]interface{}, error)
}

type pgxAgendaRepository struct {
	db *pgxpool.Pool
}

func NewAgendaRepository(db *pgxpool.Pool) AgendaRepository {
	return &pgxAgendaRepository{db: db}
}

func (r *pgxAgendaRepository) GenerateNomorAgenda(ctx context.Context, schoolID, jenis string, tahun int) (string, int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", 0, err
	}
	defer tx.Rollback(ctx)

	var lastNumber int
	queryLock := `
		SELECT last_number FROM agenda_sequences 
		WHERE school_id = $1 AND jenis = $2 AND tahun = $3 
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, queryLock, schoolID, jenis, tahun).Scan(&lastNumber)
	
	if err != nil {
		lastNumber = 0
		queryInsert := `
			INSERT INTO agenda_sequences (school_id, jenis, tahun, last_number)
			VALUES ($1, $2, $3, 0)
			ON CONFLICT (school_id, jenis, tahun) DO NOTHING
		`
		_, _ = tx.Exec(ctx, queryInsert, schoolID, jenis, tahun)
	}

	nextNumber := lastNumber + 1

	queryUpdate := `
		UPDATE agenda_sequences 
		SET last_number = $1 
		WHERE school_id = $2 AND jenis = $3 AND tahun = $4
	`
	_, err = tx.Exec(ctx, queryUpdate, nextNumber, schoolID, jenis, tahun)
	if err != nil {
		return "", 0, err
	}

	bulanRomawi := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}[int(time.Now().Month())]
	nomorAgenda := fmt.Sprintf("%03d/AGD-%s/%s/%d", nextNumber, jenis, bulanRomawi, tahun)

	if err := tx.Commit(ctx); err != nil {
		return "", 0, err
	}

	return nomorAgenda, nextNumber, nil
}

func (r *pgxAgendaRepository) CreateAgendaEntry(ctx context.Context, entry *models.AgendaEntry) error {
	query := `
		INSERT INTO agenda_entries (
			id, school_id, correspondence_id, jenis, nomor_agenda, nomor_urut, tahun, tanggal_catat, petugas_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	entry.TanggalCatat = time.Now()

	_, err := r.db.Exec(ctx, query,
		entry.ID, entry.SchoolID, entry.CorrespondenceID, entry.Jenis, entry.NomorAgenda, entry.NomorUrut, entry.Tahun, entry.TanggalCatat, entry.PetugasID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxAgendaRepository) GetByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) (*models.AgendaEntry, error) {
	query := `
		SELECT id, school_id, correspondence_id, jenis, nomor_agenda, nomor_urut, tahun, tanggal_catat, petugas_id
		FROM agenda_entries
		WHERE school_id = $1 AND correspondence_id = $2
	`
	var entry models.AgendaEntry
	err := r.db.QueryRow(ctx, query, schoolID, correspondenceID).Scan(
		&entry.ID, &entry.SchoolID, &entry.CorrespondenceID, &entry.Jenis, &entry.NomorAgenda, &entry.NomorUrut, &entry.Tahun, &entry.TanggalCatat, &entry.PetugasID,
	)
	if err != nil {
		return nil, nil
	}
	return &entry, nil
}

func (r *pgxAgendaRepository) ListAgenda(ctx context.Context, filter AgendaFilter) ([]map[string]interface{}, error) {
	// Menggunakan COALESCE untuk mencegah error scan NULL
	query := `
		SELECT 
			a.nomor_agenda, a.tanggal_catat, a.jenis,
			COALESCE(c.number, '') as number, 
			c.date, 
			COALESCE(c.subject, '') as subject, 
			COALESCE(c.from_party, '') as from_party, 
			COALESCE(c.to_party, '') as to_party, 
			COALESCE(c.kode_klasifikasi, '') as kode_klasifikasi, 
			COALESCE(c.status, '') as status
		FROM agenda_entries a
		JOIN correspondences c ON a.correspondence_id = c.id
		WHERE a.school_id = $1 AND a.tahun = $2
	`
	args := []interface{}{filter.SchoolID, filter.Tahun}
	argIndex := 3

	if filter.Jenis != "" && filter.Jenis != "BOTH" {
		query += fmt.Sprintf(" AND a.jenis = $%d", argIndex)
		args = append(args, filter.Jenis)
		argIndex++
	}

	if filter.Bulan > 0 {
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM a.tanggal_catat) = $%d", argIndex)
		args = append(args, filter.Bulan)
		argIndex++
	}

	query += " ORDER BY a.jenis, a.nomor_urut ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var (
			noAgenda, jenis, number, subject, fromParty, toParty, kodeKlasifikasi, status string
			tglCatat, tglSurat time.Time
		)
		if err := rows.Scan(&noAgenda, &tglCatat, &jenis, &number, &tglSurat, &subject, &fromParty, &toParty, &kodeKlasifikasi, &status); err != nil {
			return nil, err
		}
		
		list = append(list, map[string]interface{}{
			"nomor_agenda":     noAgenda,
			"tanggal_catat":    tglCatat.Format("2006-01-02"),
			"jenis":            jenis,
			"number":           number,
			"date":             tglSurat.Format("2006-01-02"),
			"subject":          subject,
			"from_party":       fromParty,
			"to_party":         toParty,
			"kode_klasifikasi": kodeKlasifikasi,
			"status":           status,
		})
	}
	return list, rows.Err()
}
