package errors

import "errors"

var (
	ErrInvalidAssetName       = errors.New("nama aset tidak boleh kosong")
	ErrInvalidSerialNumber    = errors.New("nomor seri tidak boleh kosong")
	ErrInvalidQRCode          = errors.New("kode QR tidak valid")
	ErrInvalidLocation        = errors.New("lokasi tidak valid (gedung/lantai/ruangan harus diisi)")
	ErrInvalidConditionChange = errors.New("perubahan kondisi tidak sah")
	ErrAssetNotFound          = errors.New("aset tidak ditemukan")
	ErrAssetAlreadyDisposed   = errors.New("aset sudah dibuang")
	ErrCannotBorrowDisposed   = errors.New("tidak dapat meminjam aset yang sudah dibuang")
	ErrCannotBorrowDamaged    = errors.New("tidak dapat meminjam aset dalam kondisi rusak")
	ErrInvalidIssueStatus     = errors.New("status issue tidak valid untuk transisi ini")
	ErrIssueAlreadyClosed     = errors.New("issue sudah ditutup")
	ErrEmptyIssueDescription  = errors.New("deskripsi issue harus diisi")
	ErrInvalidReporter        = errors.New("reporter tidak valid")
)
