package rest

import (
	"net/http"
	"strconv"
)

// registerFinanceRoutes mendaftarkan seluruh sub-router modul Keuangan
// di bawah prefix /api/v1/keuangan/*.
func (r *Router) registerFinanceRoutes(mux *http.ServeMux) {
	r.registerKeuanganMasterJenisRoutes(mux)
	r.registerKeuanganMasterCOARoutes(mux)
	r.registerKeuanganTransaksiRoutes(mux)
	r.registerKeuanganPengeluaranRoutes(mux)
	r.registerKeuanganRabRoutes(mux)
	r.registerKeuanganRekonsiliasiRoutes(mux)
	r.registerKeuanganPenutupanRoutes(mux)
	r.registerKeuanganPiutangRoutes(mux)
	r.registerKeuanganTTDRoutes(mux)
	r.registerKeuanganLaporanRoutes(mux)
}

// =====================================================================
// Shared helpers untuk router Finance/Keuangan
// =====================================================================

// resolveSchoolID mengambil school_id dari query param ?school_id= atau
// header X-School-ID. Mengembalikan string kosong bila tidak ditemukan;
// caller wajib memvalidasi.
func resolveSchoolID(req *http.Request) string {
	if v := req.URL.Query().Get("school_id"); v != "" {
		return v
	}
	if v := req.Header.Get("X-School-ID"); v != "" {
		return v
	}
	return ""
}

// parseIntQuery membaca query param int dengan fallback default bila
// kosong atau tidak valid.
func parseIntQuery(req *http.Request, key string, defaultVal int) int {
	v := req.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

// parseBoolPtrQuery membaca query param bool menjadi *bool.
// Mengembalikan nil bila kosong/invalid (artinya "tidak difilter").
func parseBoolPtrQuery(req *http.Request, key string) *bool {
	v := req.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return nil
	}
	return &b
}
