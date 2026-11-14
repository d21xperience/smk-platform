package service

import (
	"encoding/json"
	"math"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// autoMatchMutasi — algorithm auto-match mutasi bank ↔ BKU
//
// Aturan match:
//   1. |mutasi.nominal| == |bku.debet - bku.kredit|  (abs)
//   2. |selisih hari| ≤ RekonMatchToleranceHari
//   3. BKU belum terpakai
//   4. Kalau mutasi.referensi != "" → harus match bku.no_bukti
//
// Output:
//   - matched: pasangan (mutasi, bku)
//   - unmatchedBank: mutasi tidak ter-match
//   - unmatchedBuku: bku sisa tidak ter-match
// =====================================================================
func autoMatchMutasi(
	mutasi []models.MutasiParsed,
	bku []models.BkuForMatch,
) (matched []models.MatchedPair,
	unmatchedBank []models.UnmatchedBankItem,
	unmatchedBuku []models.UnmatchedBukuItem) {

	usedBku := make(map[int]bool)
	matched = []models.MatchedPair{}
	unmatchedBank = []models.UnmatchedBankItem{}

	for _, m := range mutasi {
		mutasiNominal := math.Abs(m.Debet - m.Kredit)
		found := false

		for i, b := range bku {
			if usedBku[i] {
				continue
			}
			bkuNominal := math.Abs(b.Debet - b.Kredit)
			if math.Abs(mutasiNominal-bkuNominal) > 0.01 {
				continue
			}
			selisihHari := int(math.Abs(m.Tanggal.Sub(b.Tanggal).Hours() / 24))
			if selisihHari > models.RekonMatchToleranceHari {
				continue
			}
			// Kalau mutasi punya referensi, harus match no_bukti
			if m.Referensi != "" && m.Referensi != b.NoBukti {
				continue
			}
			// Match!
			usedBku[i] = true
			matched = append(matched, models.MatchedPair{
				Mutasi:      m,
				BKU:         b,
				SelisihHari: selisihHari,
			})
			found = true
			break
		}

		if !found {
			unmatchedBank = append(unmatchedBank, models.UnmatchedBankItem{Mutasi: m})
		}
	}

	// BKU yang belum terpakai = unmatched buku
	unmatchedBuku = []models.UnmatchedBukuItem{}
	for i, b := range bku {
		if !usedBku[i] {
			unmatchedBuku = append(unmatchedBuku, models.UnmatchedBukuItem{BKU: b})
		}
	}

	return matched, unmatchedBank, unmatchedBuku
}

// =====================================================================
// computeSaldoRekonsiliasi — hitung saldo adjust kedua sisi
//
// Konvensi:
//   Sisi BANK adjust dari unmatched BUKU (item buku yang belum masuk rekening)
//   Sisi BUKU adjust dari unmatched BANK (item bank yang belum tercatat di buku)
//
// Jadi kalau di buku ada pengeluaran 1.2jt yang belum muncul di bank:
//   - saldo rekonsiliasi BANK = saldoBank + Σ(unmatched_buku.debet - kredit)
//   - saldo rekonsiliasi BUKU = saldoBuku + Σ(unmatched_bank.debet - kredit)
//
// Balance jika kedua sisi sama (dalam toleransi sen).
// =====================================================================
func computeSaldoRekonsiliasi(
	saldoBuku, saldoBank float64,
	unmatchedBank []models.UnmatchedBankItem,
	unmatchedBuku []models.UnmatchedBukuItem,
) (saldoRekBank, saldoRekBuku float64, isBalance bool) {

	// Adjust sisi BANK dengan item BUKU yang belum di bank
	saldoRekBank = saldoBank
	for _, u := range unmatchedBuku {
		saldoRekBank += u.BKU.Debet - u.BKU.Kredit
	}

	// Adjust sisi BUKU dengan item BANK yang belum di buku
	saldoRekBuku = saldoBuku
	for _, u := range unmatchedBank {
		saldoRekBuku += u.Mutasi.Debet - u.Mutasi.Kredit
	}

	// Toleransi 0.01 untuk floating-point
	isBalance = math.Abs(saldoRekBank-saldoRekBuku) < 0.01
	return saldoRekBank, saldoRekBuku, isBalance
}

// =====================================================================
// buildRekonsiliasiDetails — convert matched pairs ke RekonsiliasiDetail
// =====================================================================
func buildRekonsiliasiDetails(matched []models.MatchedPair) []models.RekonsiliasiDetail {
	out := make([]models.RekonsiliasiDetail, 0, len(matched))
	for _, m := range matched {
		out = append(out, models.RekonsiliasiDetail{
			MutasiTanggal:    m.Mutasi.Tanggal.Format("2006-01-02"),
			MutasiKeterangan: m.Mutasi.Keterangan,
			MutasiDebet:      m.Mutasi.Debet,
			MutasiKredit:     m.Mutasi.Kredit,
			MutasiSaldo:      m.Mutasi.Saldo,
			MutasiReferensi:  m.Mutasi.Referensi,
			BukuKasUmumID:    m.BKU.ID,
			SelisihHari:      m.SelisihHari,
		})
	}
	return out
}

// =====================================================================
// buildRekonsiliasiUnmatched — convert unmatched ke RekonsiliasiUnmatched
// data_json berisi snapshot data untuk audit trail
// =====================================================================
func buildRekonsiliasiUnmatched(
	unmatchedBank []models.UnmatchedBankItem,
	unmatchedBuku []models.UnmatchedBukuItem,
) []models.RekonsiliasiUnmatched {

	out := []models.RekonsiliasiUnmatched{}

	// Bank items
	for _, u := range unmatchedBank {
		payload := map[string]interface{}{
			"tanggal":    u.Mutasi.Tanggal.Format("2006-01-02"),
			"keterangan": u.Mutasi.Keterangan,
			"debet":      u.Mutasi.Debet,
			"kredit":     u.Mutasi.Kredit,
			"saldo":      u.Mutasi.Saldo,
			"referensi":  u.Mutasi.Referensi,
		}
		js, _ := json.Marshal(payload)
		out = append(out, models.RekonsiliasiUnmatched{
			Tipe:               models.UnmatchedTipeBank,
			DataJSON:           string(js),
			StatusTindakLanjut: models.UnmatchedPending,
		})
	}

	// Buku items
	for _, u := range unmatchedBuku {
		payload := map[string]interface{}{
			"tanggal":         u.BKU.Tanggal.Format("2006-01-02"),
			"no_bukti":        u.BKU.NoBukti,
			"uraian":          u.BKU.Uraian,
			"debet":           u.BKU.Debet,
			"kredit":          u.BKU.Kredit,
			"saldo_berjalan":  u.BKU.SaldoBerjalan,
			"kode_rekening_id": u.BKU.KodeRekeningID,
		}
		js, _ := json.Marshal(payload)
		out = append(out, models.RekonsiliasiUnmatched{
			Tipe:               models.UnmatchedTipeBuku,
			DataJSON:           string(js),
			StatusTindakLanjut: models.UnmatchedPending,
		})
	}

	return out
}
