package value_object

import (
	domainErrors "sekolah-platform/services/tu-core/internal/domain/inventory/errors"
	"time"
)

type AssetIdentity struct {
	ID           string
	QRCode       string
	Name         string
	SerialNumber string
	PurchaseDate time.Time
}

func NewAssetIdentity(id, qrCode, name, serial string, purchaseDate time.Time) (AssetIdentity, error) {
	if name == "" {
		return AssetIdentity{}, domainErrors.ErrInvalidAssetName
	}
	if serial == "" {
		return AssetIdentity{}, domainErrors.ErrInvalidSerialNumber
	}
	// QR Code minimal 6 karakter (validasi sederhana)
	if len(qrCode) < 6 {
		return AssetIdentity{}, domainErrors.ErrInvalidQRCode
	}
	return AssetIdentity{
		ID:           id,
		QRCode:       qrCode,
		Name:         name,
		SerialNumber: serial,
		PurchaseDate: purchaseDate,
	}, nil
}

// MatchQRCode untuk fleksibilitas identifikasi (bisa partial match)
func (i AssetIdentity) MatchQRCode(input string) bool {
	// bisa diperluas dengan regex atau fuzzy, saat ini exact match
	return i.QRCode == input
}
