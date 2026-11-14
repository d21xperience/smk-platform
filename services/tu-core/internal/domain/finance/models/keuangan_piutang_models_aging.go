package models

import "time"

// =====================================================================
// ComputeHariTerlambat — hari sejak jatuh_tempo
//
//   - Hasil positif: sudah lewat jatuh_tempo (terlambat)
//   - Hasil negatif: belum sampai jatuh_tempo (belum terlambat)
//   - Hasil 0:     jatuh_tempo hari ini
//
// Asumsi: tgl dinormalisasi ke 00:00 UTC (both input) sebelum subtract
// =====================================================================
func ComputeHariTerlambat(tanggalJatuhTempo, asOf time.Time) int {
	due := normalizeDate(tanggalJatuhTempo)
	now := normalizeDate(asOf)
	return int(now.Sub(due).Hours() / 24)
}

// =====================================================================
// ComputeAgingGroup — bucket string berdasar hari terlambat
//
//   <0       → BELUM_JATUH_TEMPO
//   0-29     → <30
//   30-59    → 30-60
//   60-89    → 60-90
//   ≥90      → >90
//
// Catatan: hari=0 masuk <30 (jatuh tempo hari ini dianggap "hampir terlambat"
// tapi masih dalam bucket pertama). Kalau bisnis ingin hari=0 masuk
// BELUM_JATUH_TEMPO, ubah kondisi pertama jadi `hari <= 0`.
// =====================================================================
func ComputeAgingGroup(tanggalJatuhTempo, asOf time.Time) string {
	hari := ComputeHariTerlambat(tanggalJatuhTempo, asOf)
	switch {
	case hari < 0:
		return PiutangAgingBelumJatuhTempo
	case hari < 30:
		return PiutangAgingLessThan30
	case hari < 60:
		return PiutangAging30to60
	case hari < 90:
		return PiutangAging60to90
	default:
		return PiutangAgingMoreThan90
	}
}

// =====================================================================
// AgingGroupPriority — untuk memilih bucket terparah (untuk OldestAgingGroup)
//
// Semakin tinggi = semakin buruk. `>90` paling buruk.
// =====================================================================
func AgingGroupPriority(group string) int {
	switch group {
	case PiutangAgingMoreThan90:
		return 5
	case PiutangAging60to90:
		return 4
	case PiutangAging30to60:
		return 3
	case PiutangAgingLessThan30:
		return 2
	case PiutangAgingBelumJatuhTempo:
		return 1
	}
	return 0
}

// normalizeDate — normalisasi time ke 00:00 UTC untuk kalkulasi hari konsisten
func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
