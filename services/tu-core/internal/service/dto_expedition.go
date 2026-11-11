package service

import "time"

type CreateExpeditionRequest struct {
	PenerimaNama     string `json:"penerima_nama"`
	PenerimaInstansi string `json:"penerima_instansi"`
	PenerimaAlamat   string `json:"penerima_alamat"`
	MetodeKirim      string `json:"metode_kirim"`
	NomorResi        string `json:"nomor_resi"`
	TanggalKirim     string `json:"tanggal_kirim"`
	JumlahLembar     int    `json:"jumlah_lembar"`
}

type ExpeditionResponse struct {
	ID               string    `json:"id"`
	CorrespondenceID string    `json:"correspondence_id"`
	PenerimaNama     string    `json:"penerima_nama"`
	PenerimaInstansi string    `json:"penerima_instansi"`
	MetodeKirim      string    `json:"metode_kirim"`
	NomorResi        string    `json:"nomor_resi"`
	TanggalKirim     time.Time `json:"tanggal_kirim"`
	JumlahLembar     int       `json:"jumlah_lembar"`
	Status           string    `json:"status"`
	BuktiPdfURL      string    `json:"bukti_pdf_url"`
	VerifyToken      string    `json:"verify_token"`
}
