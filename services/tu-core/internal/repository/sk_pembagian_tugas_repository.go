package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SKDetailAssignment struct {
	GuruName      string `json:"guruName"`
	GuruNIP       string `json:"guruNip"`
	MataPelajaran string `json:"mataPelajaran"`
	Kelas         string `json:"kelas"`
	JamPerMinggu  int    `json:"jamPerMinggu"`
}

type SKDetailTugasTambahan struct {
	GuruName   string `json:"guruName"`
	JenisTugas string `json:"jenisTugas"`
	Keterangan string `json:"keterangan"`
}

type SKDetail struct {
	ID            string                  `json:"id"`
	NomorSK       string                  `json:"nomorSk"`
	TahunAjaran   string                  `json:"tahunAjaran"`
	Semester      int                     `json:"semester"`
	Status        string                  `json:"status"`
	CreatedAt     time.Time               `json:"createdAt"`
	Assignments   []SKDetailAssignment    `json:"assignments"`
	TugasTambahan []SKDetailTugasTambahan `json:"tugasTambahan"`
}

type SKHeader struct {
	ID          string    `json:"id"`
	NomorSK     string    `json:"nomorSk"`
	TahunAjaran string    `json:"tahunAjaran"`
	Semester    int       `json:"semester"`
	Status      string    `json:"status"`
	DraftPdfURL string    `json:"draftPdfUrl"`
	FinalPdfURL string    `json:"finalPdfUrl"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SKAssignment struct {
	PtkID               string   `json:"ptkId"`
	MataPelajaranID     int      `json:"mataPelajaranId"`
	RombonganBelajarIDs []string `json:"rombonganBelajarIds"`
	JamPerMinggu        int      `json:"jamPerMinggu"`
}

type SKTugasTambahan struct {
	PtkID      string `json:"ptkId"`
	JenisTugas string `json:"jenisTugas"`
	Keterangan string `json:"keterangan"`
}

type SKCreateRequest struct {
	TahunAjaran   string            `json:"tahunAjaran"`
	Semester      int               `json:"semester"`
	Assignments   []SKAssignment    `json:"assignments"`
	TugasTambahan []SKTugasTambahan `json:"tugasTambahan"`
}

type SKPembagianTugasRepository interface {
	GenerateNextSKNumber(ctx context.Context, schoolID string) (string, error)
	CreateDraftSK(ctx context.Context, req SKCreateRequest, nomorSK, userID string) (string, error)
	ListSK(ctx context.Context) ([]SKHeader, error)
	GetDetailByID(ctx context.Context, id string) (*SKDetail, error)
	UpdateStatus(ctx context.Context, id string, status string, note string) error
	UpdateDocumentURLs(ctx context.Context, id string, draftURL, finalURL, lampiranURL, qrURL string) error
}

type pgxSKPembagianTugasRepository struct {
	db *pgxpool.Pool
}

func NewSKPembagianTugasRepository(db *pgxpool.Pool) SKPembagianTugasRepository {
	return &pgxSKPembagianTugasRepository{db: db}
}

func (r *pgxSKPembagianTugasRepository) GenerateNextSKNumber(ctx context.Context, schoolID string) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	bulan := int(now.Month())
	tahun := now.Year()

	romawi := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	bulanRomawi := romawi[bulan]

	var lastSeq int
	queryCheck := `SELECT sequence FROM sk_numbers WHERE jenis = 'SK_PEMBAGIAN_TUGAS' AND tahun = $1 AND bulan = $2 ORDER BY sequence DESC LIMIT 1 FOR UPDATE`
	err = tx.QueryRow(ctx, queryCheck, tahun, bulan).Scan(&lastSeq)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			lastSeq = 0
		} else {
			return "", err
		}
	}

	nextSeq := lastSeq + 1
	kodeSekolah := "SMK-XXX"
	nomorSK := fmt.Sprintf("421.2/SK/%s/%s/%d/%03d", kodeSekolah, bulanRomawi, tahun, nextSeq)

	_, err = tx.Exec(ctx, `
		INSERT INTO sk_numbers (jenis, tahun, bulan, sequence, nomor_sk)
		VALUES ('SK_PEMBAGIAN_TUGAS', $1, $2, $3, $4)
	`, tahun, bulan, nextSeq, nomorSK)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return nomorSK, nil
}

func (r *pgxSKPembagianTugasRepository) CreateDraftSK(ctx context.Context, req SKCreateRequest, nomorSK, userID string) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	skID := uuid.New().String()

	_, err = tx.Exec(ctx, `
		INSERT INTO sk_pembagian_tugas (id, nomor_sk, tahun_ajaran, semester, status, created_by)
		VALUES ($1, $2, $3, $4, 'DRAFT', $5)
	`, skID, nomorSK, req.TahunAjaran, req.Semester, userID)
	if err != nil {
		return "", err
	}

	for _, assign := range req.Assignments {
		detailID := uuid.New().String()
		_, err = tx.Exec(ctx, `
			INSERT INTO sk_pembagian_tugas_detail (id, sk_id, ptk_id, mata_pelajaran_id)
			VALUES ($1, $2, $3, $4)
		`, detailID, skID, assign.PtkID, assign.MataPelajaranID)
		if err != nil {
			return "", err
		}

		for _, rombelID := range assign.RombonganBelajarIDs {
			_, err = tx.Exec(ctx, `
				INSERT INTO sk_pembagian_tugas_kelas (detail_id, rombongan_belajar_id, jam_per_minggu)
				VALUES ($1, $2, $3)
			`, detailID, rombelID, assign.JamPerMinggu)
			if err != nil {
				return "", err
			}
		}
	}

	for _, tugas := range req.TugasTambahan {
		_, err = tx.Exec(ctx, `
			INSERT INTO sk_pembagian_tugas_tambahan (sk_id, ptk_id, jenis_tugas, keterangan)
			VALUES ($1, $2, $3, $4)
		`, skID, tugas.PtkID, tugas.JenisTugas, tugas.Keterangan)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return skID, nil
}

func (r *pgxSKPembagianTugasRepository) ListSK(ctx context.Context) ([]SKHeader, error) {
	query := `
		SELECT id::text, nomor_sk, tahun_ajaran, semester, status,
		       COALESCE(draft_pdf_url, ''), COALESCE(final_pdf_url, ''), created_at
		FROM sk_pembagian_tugas
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SKHeader
	for rows.Next() {
		var sk SKHeader
		if err := rows.Scan(&sk.ID, &sk.NomorSK, &sk.TahunAjaran, &sk.Semester, &sk.Status, &sk.DraftPdfURL, &sk.FinalPdfURL, &sk.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, sk)
	}
	return list, nil
}

func (r *pgxSKPembagianTugasRepository) GetDetailByID(ctx context.Context, id string) (*SKDetail, error) {
	queryHeader := `
		SELECT id::text, nomor_sk, tahun_ajaran, semester, status, created_at
		FROM sk_pembagian_tugas WHERE id = $1
	`
	var sk SKDetail
	err := r.db.QueryRow(ctx, queryHeader, id).Scan(&sk.ID, &sk.NomorSK, &sk.TahunAjaran, &sk.Semester, &sk.Status, &sk.CreatedAt)
	if err != nil {
		return nil, err
	}

	queryAssign := `
		SELECT tp.nama, COALESCE(tp.nip, ''), tm.nm_mapel, tk.nm_kelas, sk.jam_per_minggu
		FROM sk_pembagian_tugas_detail sd
		JOIN tabel_ptk tp ON sd.ptk_id = tp.ptk_id
		JOIN tabel_mapel tm ON sd.mata_pelajaran_id = tm.mata_pelajaran_id
		JOIN sk_pembagian_tugas_kelas sk ON sd.id = sk.detail_id
		JOIN tabel_kelas tk ON sk.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE sd.sk_id = $1
	`
	rows, _ := r.db.Query(ctx, queryAssign, id)
	defer rows.Close()
	for rows.Next() {
		var a SKDetailAssignment
		rows.Scan(&a.GuruName, &a.GuruNIP, &a.MataPelajaran, &a.Kelas, &a.JamPerMinggu)
		sk.Assignments = append(sk.Assignments, a)
	}

	queryTugas := `
		SELECT tp.nama, st.jenis_tugas, COALESCE(st.keterangan, '')
		FROM sk_pembagian_tugas_tambahan st
		JOIN tabel_ptk tp ON st.ptk_id = tp.ptk_id
		WHERE st.sk_id = $1
	`
	rowsT, _ := r.db.Query(ctx, queryTugas, id)
	defer rowsT.Close()
	for rowsT.Next() {
		var t SKDetailTugasTambahan
		rowsT.Scan(&t.GuruName, &t.JenisTugas, &t.Keterangan)
		sk.TugasTambahan = append(sk.TugasTambahan, t)
	}

	return &sk, nil
}

func (r *pgxSKPembagianTugasRepository) UpdateStatus(ctx context.Context, id string, status string, note string) error {
	var query string
	switch status {
	case "REJECTED":
		query = `UPDATE sk_pembagian_tugas SET status = $1, rejection_reason = $2, updated_at = NOW() WHERE id = $3`
	case "REVISION":
		query = `UPDATE sk_pembagian_tugas SET status = $1, revision_notes = $2, updated_at = NOW() WHERE id = $3`
	default:
		query = `UPDATE sk_pembagian_tugas SET status = $1, updated_at = NOW() WHERE id = $2`
		_, err := r.db.Exec(ctx, query, status, id)
		return err
	}
	_, err := r.db.Exec(ctx, query, status, note, id)
	return err
}

func (r *pgxSKPembagianTugasRepository) UpdateDocumentURLs(ctx context.Context, id string, draftURL, finalURL, lampiranURL, qrURL string) error {
	query := `
		UPDATE sk_pembagian_tugas
		SET draft_pdf_url = COALESCE(NULLIF($1, ''), draft_pdf_url),
		    final_pdf_url = COALESCE(NULLIF($2, ''), final_pdf_url),
		    lampiran_pdf_url = COALESCE(NULLIF($3, ''), lampiran_pdf_url),
		    qr_tte_url = COALESCE(NULLIF($4, ''), qr_tte_url),
		    updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, draftURL, finalURL, lampiranURL, qrURL, id)
	return err
}
