package models

import "time"

type Expedition struct {
	ID               string
	SchoolID         string
	CorrespondenceID string
	PenerimaNama     string
	PenerimaInstansi string
	PenerimaAlamat   string
	MetodeKirim      string
	NomorResi        string
	TanggalKirim     time.Time
	JumlahLembar     int
	Status           string
	BuktiPdfURL      string
	VerifyToken      string
	CreatedBy        string
	CreatedAt        time.Time
}
