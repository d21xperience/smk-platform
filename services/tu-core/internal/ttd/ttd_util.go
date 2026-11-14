// Package ttd menyediakan util untuk tanda tangan digital (TTD):
// generation hash SHA-256 dan pembacaan secret dari environment.
//
// Fungsi di package ini pure (tidak ada dependency ke service/repo),
// sehingga bisa dipakai lintas modul (penutupan, pengeluaran, dll).
package ttd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// env var untuk secret
	ttdSecretEnv = "TTD_SECRET"

	// fallback dev — HANYA untuk development. Di production user WAJIB set TTD_SECRET.
	ttdSecretFallback = "smk-platform-dev-secret-change-me"
)

var (
	secretOnce   sync.Once
	cachedSecret string
)

// GetTTDSecret — baca ENV TTD_SECRET sekali, fallback ke hardcoded dev secret.
// Warning di-log sekali saat fallback digunakan, agar developer sadar.
func GetTTDSecret() string {
	secretOnce.Do(func() {
		if s := os.Getenv(ttdSecretEnv); s != "" {
			cachedSecret = s
			return
		}
		cachedSecret = ttdSecretFallback
		log.Printf("⚠️  WARNING: ENV %s tidak diset, pakai dev fallback. JANGAN dipakai di production!",
			ttdSecretEnv)
	})
	return cachedSecret
}

// GenerateSignature — SHA-256 hex dari kombinasi (docType, docID, userID, timestamp, secret).
//
//	docType   : string pendek, misal "PENUTUPAN", "PENGELUARAN"
//	docID     : UUID dokumen
//	userID    : UUID signer (kepala sekolah / pejabat)
//	timestamp : waktu sign (mikrodetik presisi)
//
// Output: 64-character hex string (SHA-256 raw bytes di-encode hex).
func GenerateSignature(docType, docID, userID string, ts time.Time) string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%s",
		docType,
		docID,
		userID,
		ts.UTC().Format(time.RFC3339Nano),
		GetTTDSecret(),
	)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

// VerifySignature — verifikasi hash deterministic.
// Return true kalau hash input == GenerateSignature(params yang sama).
//
// Catatan: untuk endpoint publik verify-by-hash, kita TIDAK pakai ini
// (kita lookup langsung ke kolom ttd_kepsek_hash di DB). Fungsi ini
// berguna untuk unit test atau validator internal.
func VerifySignature(docType, docID, userID string, ts time.Time, expectedHash string) bool {
	return GenerateSignature(docType, docID, userID, ts) == expectedHash
}
