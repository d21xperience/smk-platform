package service

import "time"

type LetterTypeRequest struct {
	Kategori        string `json:"kategori"`
	Peruntukan      string `json:"peruntukan"`
	Value           string `json:"value"`
	Label           string `json:"label"`
	KodeKlasifikasi string `json:"kode_klasifikasi"`
}

type LetterTypeResponse struct {
	ID              string    `json:"id"`
	Kategori        string    `json:"kategori"`
	Peruntukan      string    `json:"peruntukan"`
	Value           string    `json:"value"`
	Label           string    `json:"label"`
	KodeKlasifikasi string    `json:"kode_klasifikasi"`
	CreatedAt       time.Time `json:"created_at"`
}
