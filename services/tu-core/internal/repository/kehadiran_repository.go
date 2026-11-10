package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KehadiranHarian struct {
	ID             string  `json:"id"`
	PtkID          string  `json:"ptkId"`
	NamaPegawai    string  `json:"namaPegawai"`
	NIP            string  `json:"nip"`
	JenisGTK       string  `json:"jenisGtk"`
	Tanggal        string  `json:"tanggal"`
	Status         string  `json:"status"`
	JamMasuk       *string `json:"jamMasuk"`
	JamPulang      *string `json:"jamPulang"`
	DurasiKerja    *int    `json:"durasiKerja"`
	MenitTerlambat *int    `json:"menitTerlambat"`
	SourceData     string  `json:"sourceData"`
	IsCorrected    bool    `json:"isCorrected"`
}

type RekapBulananHeader struct {
	ID              string  `json:"id"`
	Bulan           int     `json:"bulan"`
	Tahun           int     `json:"tahun"`
	TotalPegawai    int     `json:"totalPegawai"`
	RataPersentase  float64 `json:"rataPersentase"`
	Status          string  `json:"status"`
	ExcelURL        string  `json:"excelUrl"`
}

type KehadiranRepository interface {
	GetAttendanceByDate(ctx context.Context, tanggal time.Time) ([]KehadiranHarian, error)
	GetActiveEmployees(ctx context.Context) ([]map[string]interface{}, error)
	GetApprovedLeavesForDate(ctx context.Context, tanggal time.Time) (map[string]string, error)
	UpsertAttendance(ctx context.Context, ptkID string, tanggal time.Time, status string, jamMasuk, jamPulang *string, durasi, menitTerlambat *int, source string) error
	AutoCollectDaily(ctx context.Context, tanggal time.Time) (int, error)
	GetRekapBulanan(ctx context.Context) ([]RekapBulananHeader, error)
	GenerateRekapBulanan(ctx context.Context, bulan, tahun int) (string, error)
}

type pgxKehadiranRepository struct {
	db *pgxpool.Pool
}

func NewKehadiranRepository(db *pgxpool.Pool) KehadiranRepository {
	return &pgxKehadiranRepository{db: db}
}

func (r *pgxKehadiranRepository) GetAttendanceByDate(ctx context.Context, tanggal time.Time) ([]KehadiranHarian, error) {
	query := `
		SELECT 
			COALESCE(pa.id::text, gen_random_uuid()::text),
			tp.ptk_id::text,
			COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, ''),
			COALESCE(tp.nip, ''),
			CASE WHEN tp.jenis_ptk_id = 1 THEN 'GURU' ELSE 'TENDIK' END,
			to_char($1::date, 'YYYY-MM-DD'),
			COALESCE(pa.status, 'BELUM_ABSEN'),
			to_char(pa.jam_masuk, 'HH24:MI'),
			to_char(pa.jam_pulang, 'HH24:MI'),
			pa.durasi_kerja, pa.menit_terlambat,
			COALESCE(pa.source_data, '-'),
			COALESCE(pa.is_corrected, false)
		FROM tabel_ptk tp
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		LEFT JOIN ptk_attendance pa ON tp.ptk_id = pa.ptk_id AND pa.date = $1
		WHERE tp.soft_delete = 0
		ORDER BY tp.nama ASC
	`
	rows, err := r.db.Query(ctx, query, tanggal)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]KehadiranHarian, 0)
	for rows.Next() {
		var k KehadiranHarian
		if err := rows.Scan(
			&k.ID, &k.PtkID, &k.NamaPegawai, &k.NIP, &k.JenisGTK,
			&k.Tanggal, &k.Status, &k.JamMasuk, &k.JamPulang,
			&k.DurasiKerja, &k.MenitTerlambat, &k.SourceData, &k.IsCorrected,
		); err != nil {
			return nil, err
		}
		list = append(list, k)
	}
	return list, nil
}

func (r *pgxKehadiranRepository) GetActiveEmployees(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT tp.ptk_id::text, 
		       CASE WHEN tp.jenis_ptk_id = 1 THEN 'GURU' ELSE 'TENDIK' END
		FROM tabel_ptk tp
		WHERE tp.soft_delete = 0
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {
		var ptkID, jenis string
		if err := rows.Scan(&ptkID, &jenis); err == nil {
			list = append(list, map[string]interface{}{"ptkId": ptkID, "jenis": jenis})
		}
	}
	return list, nil
}

func (r *pgxKehadiranRepository) GetApprovedLeavesForDate(ctx context.Context, tanggal time.Time) (map[string]string, error) {
	query := `
		SELECT ptk_id::text, jenis_cuti 
		FROM ptk_leaves 
		WHERE status = 'disetujui' 
		  AND tanggal_mulai <= $1 AND tanggal_selesai >= $1
	`
	rows, err := r.db.Query(ctx, query, tanggal)
	if err != nil { return nil, err }
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var ptkID, jenis string
		if err := rows.Scan(&ptkID, &jenis); err == nil {
			status := "IZIN_RESMI"
			if jenis == "Cuti Sakit" { status = "SAKIT" }
			if jenis == "SPT (Surat Perintah Tugas)" { status = "DINAS_LUAR" }
			result[ptkID] = status
		}
	}
	return result, nil
}

func (r *pgxKehadiranRepository) UpsertAttendance(ctx context.Context, ptkID string, tanggal time.Time, status string, jamMasuk, jamPulang *string, durasi, menitTerlambat *int, source string) error {
	query := `
		INSERT INTO ptk_attendance (id, ptk_id, date, status, jam_masuk, jam_pulang, durasi_kerja, menit_terlambat, source_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (ptk_id, date) DO UPDATE SET
			status = EXCLUDED.status, jam_masuk = EXCLUDED.jam_masuk, jam_pulang = EXCLUDED.jam_pulang,
			durasi_kerja = EXCLUDED.durasi_kerja, menit_terlambat = EXCLUDED.menit_terlambat,
			source_data = EXCLUDED.source_data
	`
	_, err := r.db.Exec(ctx, query, uuid.New().String(), ptkID, tanggal, status, jamMasuk, jamPulang, durasi, menitTerlambat, source)
	return err
}

func (r *pgxKehadiranRepository) AutoCollectDaily(ctx context.Context, tanggal time.Time) (int, error) {
	// 1. Ambil pegawai aktif
	pegawaiList, err := r.GetActiveEmployees(ctx)
	if err != nil { return 0, err }

	// 2. Ambil cuti yang disetujui untuk tanggal ini
	leaves, _ := r.GetApprovedLeavesForDate(ctx, tanggal)

	count := 0
	for _, p := range pegawaiList {
		ptkID := p["ptkId"].(string)

		// Cek apakah ada cuti
		if status, ok := leaves[ptkID]; ok {
			if err := r.UpsertAttendance(ctx, ptkID, tanggal, status, nil, nil, nil, nil, "CUTI_APPROVED"); err == nil {
				count++
			}
			continue
		}

		// Default: ALPHA (karena belum ada integrasi fingerprint)
		if err := r.UpsertAttendance(ctx, ptkID, tanggal, "ALPHA", nil, nil, nil, nil, "AUTO_COLLECT"); err == nil {
			count++
		}
	}

	return count, nil
}

func (r *pgxKehadiranRepository) GetRekapBulanan(ctx context.Context) ([]RekapBulananHeader, error) {
	query := `
		SELECT id::text, bulan, tahun, COALESCE(total_pegawai, 0), 
		       COALESCE(rata_persentase, 0), status, COALESCE(excel_url, '')
		FROM ptk_rekap_bulanan
		ORDER BY tahun DESC, bulan DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	list := make([]RekapBulananHeader, 0)
	for rows.Next() {
		var h RekapBulananHeader
		if err := rows.Scan(&h.ID, &h.Bulan, &h.Tahun, &h.TotalPegawai, &h.RataPersentase, &h.Status, &h.ExcelURL); err == nil {
			list = append(list, h)
		}
	}
	return list, nil
}

func (r *pgxKehadiranRepository) GenerateRekapBulanan(ctx context.Context, bulan, tahun int) (string, error) {
	rekapID := uuid.New().String()

	// 1. Insert Header
	queryHeader := `
		INSERT INTO ptk_rekap_bulanan (id, bulan, tahun, status)
		VALUES ($1, $2, $3, 'GENERATING')
		ON CONFLICT (bulan, tahun) DO UPDATE SET status = 'GENERATING', generated_at = NOW()
		RETURNING id::text
	`
	err := r.db.QueryRow(ctx, queryHeader, rekapID, bulan, tahun).Scan(&rekapID)
	if err != nil { return "", err }

	// 2. Hapus detail lama jika regenerate
	_, _ = r.db.Exec(ctx, "DELETE FROM ptk_rekap_bulanan_detail WHERE rekap_id = $1", rekapID)

	// 3. Hitung agregat per pegawai
	queryDetail := `
		INSERT INTO ptk_rekap_bulanan_detail (
			id, rekap_id, ptk_id, total_hadir, total_terlambat, 
			total_sakit, total_izin, total_alpha, total_dinas_luar, persentase_kehadiran
		)
		SELECT 
			gen_random_uuid(), $1, tp.ptk_id,
			COUNT(CASE WHEN pa.status IN ('HADIR', 'TERLAMBAT') THEN 1 END) AS hadir,
			COUNT(CASE WHEN pa.status = 'TERLAMBAT' THEN 1 END) AS terlambat,
			COUNT(CASE WHEN pa.status = 'SAKIT' THEN 1 END) AS sakit,
			COUNT(CASE WHEN pa.status IN ('IZIN_RESMI', 'IZIN_PRIBADI') THEN 1 END) AS izin,
			COUNT(CASE WHEN pa.status = 'ALPHA' THEN 1 END) AS alpha,
			COUNT(CASE WHEN pa.status = 'DINAS_LUAR' THEN 1 END) AS dinas_luar,
			CASE 
				WHEN COUNT(pa.id) > 0 THEN 
					ROUND((COUNT(CASE WHEN pa.status IN ('HADIR', 'TERLAMBAT', 'DINAS_LUAR') THEN 1 END)::numeric / COUNT(pa.id)::numeric) * 100, 2)
				ELSE 0 
			END AS persentase
		FROM tabel_ptk tp
		LEFT JOIN ptk_attendance pa ON tp.ptk_id = pa.ptk_id 
			AND EXTRACT(MONTH FROM pa.date) = $2 
			AND EXTRACT(YEAR FROM pa.date) = $3
		WHERE tp.soft_delete = 0
		GROUP BY tp.ptk_id
	`
	_, err = r.db.Exec(ctx, queryDetail, rekapID, bulan, tahun)
	if err != nil { return "", err }

	// 4. Update Header dengan agregat
	queryUpdate := `
		UPDATE ptk_rekap_bulanan SET
			total_pegawai = (SELECT COUNT(*) FROM ptk_rekap_bulanan_detail WHERE rekap_id = $1),
			rata_persentase = (SELECT COALESCE(AVG(persentase_kehadiran), 0) FROM ptk_rekap_bulanan_detail WHERE rekap_id = $1),
			status = 'COMPLETED'
		WHERE id = $1
	`
	_, err = r.db.Exec(ctx, queryUpdate, rekapID)

	return rekapID, err
}
