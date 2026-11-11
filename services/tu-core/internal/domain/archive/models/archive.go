package models

import "time"

type Archive struct {
	ID               string
	SchoolID         string
	CorrespondenceID string
	LokasiKode       string
	Kategori         string
	TanggalArsip     time.Time
	TanggalMusnah    time.Time
	RetensiTahun     int
	Status           string
	CurrentHolderID  string
	Keterangan       string
	ArchivedBy       string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ArchiveMovement struct {
	ID                      string
	ArchiveID               string
	Type                    string // ARCHIVE, BORROW, RETURN, MOVE, DESTROY
	PetugasID               string
	PeminjamID              string
	Tanggal                 time.Time
	TanggalKembaliEstimasi  *time.Time
	TanggalKembaliAktual    *time.Time
	Tujuan                  string
	Keterangan              string
	Status                  string // AKTIF, SELESAI
}

type RetentionSchedule struct {
	ID            string
	Kategori      string
	RetensiTahun  int
	Keterangan    string
	CreatedAt     time.Time
}
