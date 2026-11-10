package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SKPTemplate struct {
	ID                  string  `json:"id"`
	JenisPegawai        string  `json:"jenisPegawai"`
	NamaTemplate        string  `json:"namaTemplate"`
	BobotPerilaku       int     `json:"bobotPerilaku"`
	BobotHasilUtama     int     `json:"bobotHasilUtama"`
	BobotHasilTambahan  int     `json:"bobotHasilTambahan"`
}

type SKPPeriode struct {
	ID               string    `json:"id"`
	NamaPeriode      string    `json:"namaPeriode"`
	JenisPeriode     string    `json:"jenisPeriode"`
	TanggalMulai     string    `json:"tanggalMulai"`
	TanggalSelesai   string    `json:"tanggalSelesai"`
	DeadlinePengisian string    `json:"deadlinePengisian"`
	TahunReferensi   int       `json:"tahunReferensi"`
	Status           string    `json:"status"`
}

type SKPPenilaian struct {
	ID              string  `json:"id"`
	PeriodeID       string  `json:"periodeId"`
	PtkID           string  `json:"ptkId"`
	TeacherName     string  `json:"teacherName"`
	TemplateID      string  `json:"templateId"`
	PenilaiUtamaID  string  `json:"penilaiUtamaId"`
	PenilaiBerikatID string `json:"penilaiBerikatId"`
	Status          string  `json:"status"`
	NilaiAkhir      float64 `json:"nilaiAkhir"`
	PredikatAkhir   string  `json:"predikatAkhir"`
}

type SKPDetailItem struct {
	Kelompok string `json:"kelompok"`
	Aspek    string `json:"aspek"`
	Nilai    int    `json:"nilai"`
}

type SKPHasil struct {
	ID               string  `json:"id"`
	PenilaianID      string  `json:"penilaianId"`
	NilaiPerilaku    float64 `json:"nilaiPerilaku"`
	PredikatPerilaku string  `json:"predikatPerilaku"`
	NilaiHasilUtama  float64 `json:"nilaiHasilUtama"`
	NilaiHasilTambahan float64 `json:"nilaiHasilTambahan"`
	NilaiAkhir       float64 `json:"nilaiAkhir"`
	PredikatAkhir    string  `json:"predikatAkhir"`
	CatatanPenilai   string  `json:"catatanPenilai"`
	CatatanBerikat   string  `json:"catatanBerikat"`
}

type SKPRepository interface {
	ListTemplates(ctx context.Context) ([]SKPTemplate, error)
	ListPeriodes(ctx context.Context) ([]SKPPeriode, error)
	CreatePeriode(ctx context.Context, req SKPPeriode) (string, error)
	ListPenilaianByPeriode(ctx context.Context, periodeID string) ([]SKPPenilaian, error)
	GetPenilaianByID(ctx context.Context, id string) (*SKPPenilaian, error)
	GetHasilByPenilaianID(ctx context.Context, penilaianID string) (*SKPHasil, error)
	GetDetailByPenilaianID(ctx context.Context, penilaianID string) ([]SKPDetailItem, error)
	SaveHasil(ctx context.Context, hasil *SKPHasil, detail []SKPDetailItem) error
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateStatusWithCatatan(ctx context.Context, id string, status string, catatan string) error
}

type pgxSKPRepository struct {
	db *pgxpool.Pool
}

func NewSKPRepository(db *pgxpool.Pool) SKPRepository {
	return &pgxSKPRepository{db: db}
}

func (r *pgxSKPRepository) ListTemplates(ctx context.Context) ([]SKPTemplate, error) {
	query := `SELECT id::text, jenis_pegawai, nama_template, bobot_perilaku, bobot_hasil_utama, bobot_hasil_tambahan FROM penilaian_template WHERE is_active = true`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []SKPTemplate
	for rows.Next() {
		var t SKPTemplate
		if err := rows.Scan(&t.ID, &t.JenisPegawai, &t.NamaTemplate, &t.BobotPerilaku, &t.BobotHasilUtama, &t.BobotHasilTambahan); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, nil
}

func (r *pgxSKPRepository) ListPeriodes(ctx context.Context) ([]SKPPeriode, error) {
	// Menggunakan to_char untuk mengubah DATE menjadi string YYYY-MM-DD
	query := `
		SELECT id::text, nama_periode, jenis_periode, 
		       to_char(tanggal_mulai, 'YYYY-MM-DD'), 
		       to_char(tanggal_selesai, 'YYYY-MM-DD'), 
		       to_char(deadline_pengisian, 'YYYY-MM-DD'), 
		       tahun_referensi, status 
		FROM penilaian_periode 
		ORDER BY tahun_referensi DESC, created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]SKPPeriode, 0) // Inisialisasi array kosong agar tidak null
	for rows.Next() {
		var p SKPPeriode
		if err := rows.Scan(&p.ID, &p.NamaPeriode, &p.JenisPeriode, &p.TanggalMulai, &p.TanggalSelesai, &p.DeadlinePengisian, &p.TahunReferensi, &p.Status); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *pgxSKPRepository) CreatePeriode(ctx context.Context, req SKPPeriode) (string, error) {
	id := uuid.New().String()
	// PostgreSQL dapat melakukan cast otomatis dari string 'YYYY-MM-DD' ke DATE
	query := `INSERT INTO penilaian_periode (id, nama_periode, jenis_periode, tanggal_mulai, tanggal_selesai, deadline_pengisian, tahun_referensi, status) VALUES ($1, $2, $3, $4, $5, $6, $7, 'OPEN')`
	_, err := r.db.Exec(ctx, query, id, req.NamaPeriode, req.JenisPeriode, req.TanggalMulai, req.TanggalSelesai, req.DeadlinePengisian, req.TahunReferensi)
	return id, err
}

func (r *pgxSKPRepository) ListPenilaianByPeriode(ctx context.Context, periodeID string) ([]SKPPenilaian, error) {
	query := `
		SELECT pp.id::text, pp.periode_id::text, pp.ptk_id::text, 
		       COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_guru,
		       pp.template_id::text, COALESCE(pp.penilai_utama_id::text, ''), COALESCE(pp.penilai_berikat_id::text, ''), pp.status,
		       COALESCE(h.nilai_akhir, 0), COALESCE(h.predikat_akhir, '')
		FROM penilaian_pegawai pp
		JOIN tabel_ptk tp ON pp.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		LEFT JOIN penilaian_hasil h ON pp.id = h.penilaian_id
		WHERE pp.periode_id = $1
	`
	rows, err := r.db.Query(ctx, query, periodeID)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]SKPPenilaian, 0)
	for rows.Next() {
		var p SKPPenilaian
		if err := rows.Scan(&p.ID, &p.PeriodeID, &p.PtkID, &p.TeacherName, &p.TemplateID, &p.PenilaiUtamaID, &p.PenilaiBerikatID, &p.Status, &p.NilaiAkhir, &p.PredikatAkhir); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *pgxSKPRepository) GetPenilaianByID(ctx context.Context, id string) (*SKPPenilaian, error) {
	query := `
		SELECT pp.id::text, pp.periode_id::text, pp.ptk_id::text, tp.nama, pp.template_id::text, 
		       COALESCE(pp.penilai_utama_id::text, ''), COALESCE(pp.penilai_berikat_id::text, ''), pp.status
		FROM penilaian_pegawai pp
		JOIN tabel_ptk tp ON pp.ptk_id = tp.ptk_id
		WHERE pp.id = $1
	`
	var p SKPPenilaian
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.PeriodeID, &p.PtkID, &p.TeacherName, &p.TemplateID, &p.PenilaiUtamaID, &p.PenilaiBerikatID, &p.Status)
	if err != nil { return nil, err }
	return &p, nil
}

func (r *pgxSKPRepository) GetHasilByPenilaianID(ctx context.Context, penilaianID string) (*SKPHasil, error) {
	query := `
		SELECT id::text, penilaian_id::text, nilai_perilaku, predikat_perilaku, nilai_hasil_utama, COALESCE(nilai_hasil_tambahan, 0), 
		       nilai_akhir, predikat_akhir, COALESCE(catatan_penilai, ''), COALESCE(catatan_berikat, '')
		FROM penilaian_hasil WHERE penilaian_id = $1
	`
	var h SKPHasil
	err := r.db.QueryRow(ctx, query, penilaianID).Scan(&h.ID, &h.PenilaianID, &h.NilaiPerilaku, &h.PredikatPerilaku, &h.NilaiHasilUtama, &h.NilaiHasilTambahan, &h.NilaiAkhir, &h.PredikatAkhir, &h.CatatanPenilai, &h.CatatanBerikat)
	if err != nil { return nil, err }
	return &h, nil
}

func (r *pgxSKPRepository) GetDetailByPenilaianID(ctx context.Context, penilaianID string) ([]SKPDetailItem, error) {
	query := `SELECT kelompok, aspek, nilai FROM penilaian_detail WHERE penilaian_id = $1`
	rows, err := r.db.Query(ctx, query, penilaianID)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]SKPDetailItem, 0)
	for rows.Next() {
		var d SKPDetailItem
		if err := rows.Scan(&d.Kelompok, &d.Aspek, &d.Nilai); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *pgxSKPRepository) SaveHasil(ctx context.Context, hasil *SKPHasil, detail []SKPDetailItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	if hasil.ID == "" { hasil.ID = uuid.New().String() }

	query := `
		INSERT INTO penilaian_hasil (id, penilaian_id, nilai_perilaku, predikat_perilaku, nilai_hasil_utama, nilai_hasil_tambahan, nilai_akhir, predikat_akhir, catatan_penilai)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (penilaian_id) DO UPDATE SET
			nilai_perilaku = EXCLUDED.nilai_perilaku, predikat_perilaku = EXCLUDED.predikat_perilaku,
			nilai_hasil_utama = EXCLUDED.nilai_hasil_utama, nilai_hasil_tambahan = EXCLUDED.nilai_hasil_tambahan,
			nilai_akhir = EXCLUDED.nilai_akhir, predikat_akhir = EXCLUDED.predikat_akhir,
			catatan_penilai = EXCLUDED.catatan_penilai, updated_at = NOW()
	`
	_, err = tx.Exec(ctx, query, hasil.ID, hasil.PenilaianID, hasil.NilaiPerilaku, hasil.PredikatPerilaku, hasil.NilaiHasilUtama, hasil.NilaiHasilTambahan, hasil.NilaiAkhir, hasil.PredikatAkhir, hasil.CatatanPenilai)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `DELETE FROM penilaian_detail WHERE penilaian_id = $1`, hasil.PenilaianID)
	if err != nil { return err }

	for _, d := range detail {
		detailID := uuid.New().String()
		_, err = tx.Exec(ctx, `INSERT INTO penilaian_detail (id, penilaian_id, kelompok, aspek, nilai) VALUES ($1, $2, $3, $4, $5)`, detailID, hasil.PenilaianID, d.Kelompok, d.Aspek, d.Nilai)
		if err != nil { return err }
	}

	return tx.Commit(ctx)
}

func (r *pgxSKPRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE penilaian_pegawai SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *pgxSKPRepository) UpdateStatusWithCatatan(ctx context.Context, id string, status string, catatan string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE penilaian_pegawai SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `UPDATE penilaian_hasil SET catatan_berikat = $1 WHERE penilaian_id = $2`, catatan, id)
	if err != nil { return err }

	return tx.Commit(ctx)
}
