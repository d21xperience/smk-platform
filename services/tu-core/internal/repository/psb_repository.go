package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CalonSiswa struct {
	ID        string `json:"id"`
	NIS       string `json:"nis"`
	Nama      string `json:"nama"`
	Gender    string `json:"gender"`
	JurusanID string `json:"jurusan_id"`
}

type PSBRepository interface {
	GetCalonSiswa(ctx context.Context, schoolID, periodID string) ([]CalonSiswa, error)
	GenerateKelas(ctx context.Context, schoolID, periodID string, kelasData map[string][]CalonSiswa) error
}

type pgxPSBRepository struct {
	db *pgxpool.Pool
}

func NewPSBRepository(db *pgxpool.Pool) PSBRepository {
	return &pgxPSBRepository{db: db}
}

func (r *pgxPSBRepository) GetCalonSiswa(ctx context.Context, schoolID, periodID string) ([]CalonSiswa, error) {
	// Mengambil data dari tabel registrations yang sudah memiliki NIS
	query := `
		SELECT 
			id::text, 
			nis, 
			calon_nama, 
			calon_jenis_kelamin, 
			jurusan_dipilih 
		FROM registrations
		WHERE school_id = $1 
		  AND academic_period_id = $2 
		  AND nis IS NOT NULL 
		  AND status != 'ACTIVATED'
	`
	rows, err := r.db.Query(ctx, query, schoolID, periodID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var list []CalonSiswa
	for rows.Next() {
		var s CalonSiswa
		if err := rows.Scan(&s.ID, &s.NIS, &s.Nama, &s.Gender, &s.JurusanID); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *pgxPSBRepository) GenerateKelas(ctx context.Context, schoolID, periodID string, kelasData map[string][]CalonSiswa) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	for namaKelas, siswaList := range kelasData {
		if len(siswaList) == 0 {
			continue
		}

		var rombelID string
		queryCariRombel := `SELECT rombongan_belajar_id FROM tabel_kelas WHERE nm_kelas = $1 AND semester_id = $2 LIMIT 1`
		_ = tx.QueryRow(ctx, queryCariRombel, "X "+namaKelas, periodID).Scan(&rombelID)

		if rombelID == "" {
			rombelID = uuid.New().String()
			insertRombel := `
				INSERT INTO tabel_kelas (
					rombongan_belajar_id, sekolah_id, semester_id, nm_kelas, 
					tingkat_pendidikan_id, jenis_rombel, kurikulum_id, capacity
				) VALUES ($1, $2, $3, $4, 10, 1, 1, 36)
			`
			_, err = tx.Exec(ctx, insertRombel, rombelID, schoolID, periodID, "X "+namaKelas)
			if err != nil {
				return fmt.Errorf("gagal buat rombel: %w", err)
			}
		}

		insertAnggota := `
			INSERT INTO tabel_anggotakelas (anggota_rombel_id, peserta_didik_id, rombongan_belajar_id, semester_id) 
			VALUES ($1, $2, $3, $4)
		`
		for _, s := range siswaList {
			// Disini kita perlu ID Peserta Didik dari tabel_siswa. 
			// Karena data registrasi belum tentu ada di tabel_siswa, kita perlu insert atau join.
			// Asumsi: Ketika registrasi diaktifkan, siswa sudah masuk ke tabel_siswa.
			var pesertaDidikID string
			errScan := tx.QueryRow(ctx, `SELECT peserta_didik_id FROM tabel_siswa WHERE nis = $1`, s.NIS).Scan(&pesertaDidikID)
			
			if errScan == nil {
				_, err = tx.Exec(ctx, insertAnggota, uuid.New().String(), pesertaDidikID, rombelID, periodID)
				if err != nil {
					return fmt.Errorf("gagal insert anggota: %w", err)
				}
			}
		}
	}

	return tx.Commit(ctx)
}
