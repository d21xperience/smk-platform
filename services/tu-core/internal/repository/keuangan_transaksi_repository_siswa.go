package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// FindSiswaByNIS — lookup siswa dari tabel Dapodik untuk validasi NIS.
// Jika NIS tidak ditemukan, return error (bukan nil, nil) — caller tidak
// perlu cek nil sebelum cek error.
//
// Catatan: NIS bisa saja tidak ada karena siswa belum disinkron dari Dapodik.
// Kalau ini terjadi, caller (service) harus informasikan ke user agar
// menyelesaikan sinkronisasi Dapodik terlebih dahulu.
func (r *pgxKeuanganTransaksiRepository) FindSiswaByNIS(
	ctx context.Context, schoolID, nis string,
) (*models.SiswaInfo, error) {
	query := `
		SELECT
			s.peserta_didik_id::text,
			s.nis,
			s.nm_siswa,
			COALESCE(
				(SELECT k.nm_kelas
				 FROM tabel_anggotakelas ak
				 JOIN tabel_kelas k ON k.rombongan_belajar_id = ak.rombongan_belajar_id
				 WHERE ak.peserta_didik_id = s.peserta_didik_id
				 LIMIT 1),
				''
			) AS kelas_nama
		FROM tabel_siswa s
		WHERE s.nis = $1
		LIMIT 1
	`
	var info models.SiswaInfo
	err := r.db.QueryRow(ctx, query, nis).Scan(
		&info.PesertaDidikID, &info.NIS, &info.NamaSiswa, &info.KelasNama,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("NIS %s tidak ditemukan di tabel siswa (mungkin belum disinkron dari Dapodik)", nis)
	}
	if err != nil {
		return nil, fmt.Errorf("gagal lookup siswa: %w", err)
	}
	return &info, nil
}
