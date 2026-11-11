package service

import "time"

type ArchiveCorrespondenceRequest struct {
	LokasiRak       string `json:"lokasi_rak"`       // Contoh: RAK-A
	LokasiBox       string `json:"lokasi_box"`       // Contoh: BOX-01
	LokasiFolder    string `json:"lokasi_folder"`    // Contoh: FOLDER-03
	Kategori        string `json:"kategori"`         // AKTIF, INAKTIF, STATIS
	RetensiTahun    int    `json:"retensi_tahun"`
	Keterangan      string `json:"keterangan"`
}

type BorrowArchiveRequest struct {
	PeminjamID            string `json:"peminjam_id"`
	Tujuan                string `json:"tujuan"`
	TanggalKembaliEstimasi string `json:"tanggal_kembali_estimasi"` // YYYY-MM-DD
}

type MoveArchiveRequest struct {
	LokasiRak    string `json:"lokasi_rak"`
	LokasiBox    string `json:"lokasi_box"`
	LokasiFolder string `json:"lokasi_folder"`
	Keterangan   string `json:"keterangan"`
}

type ArchiveResponse struct {
	ID               string    `json:"id"`
	CorrespondenceID string    `json:"correspondence_id"`
	LokasiKode       string    `json:"lokasi_kode"`
	Kategori         string    `json:"kategori"`
	TanggalArsip     time.Time `json:"tanggal_arsip"`
	TanggalMusnah    time.Time `json:"tanggal_musnah"`
	RetensiTahun     int       `json:"retensi_tahun"`
	Status           string    `json:"status"`
	CurrentHolderID  string    `json:"current_holder_id"`
	Keterangan       string    `json:"keterangan"`
	ArchivedBy       string    `json:"archived_by"`
}

type ArchiveMovementResponse struct {
	ID                     string     `json:"id"`
	ArchiveID              string     `json:"archive_id"`
	Type                   string     `json:"type"`
	PetugasID              string     `json:"petugas_id"`
	PeminjamID             string     `json:"peminjam_id"`
	Tanggal                time.Time  `json:"tanggal"`
	TanggalKembaliEstimasi *time.Time `json:"tanggal_kembali_estimasi"`
	TanggalKembaliAktual   *time.Time `json:"tanggal_kembali_aktual"`
	Tujuan                 string     `json:"tujuan"`
	Keterangan             string     `json:"keterangan"`
	Status                 string     `json:"status"`
}
