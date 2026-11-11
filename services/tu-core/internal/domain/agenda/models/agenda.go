package models

import "time"

type AgendaEntry struct {
	ID               string
	SchoolID         string
	CorrespondenceID string
	Jenis            string // IN, OUT
	NomorAgenda      string
	NomorUrut        int
	Tahun            int
	TanggalCatat     time.Time
	PetugasID        string
}
