// Package context menyediakan Operational Context untuk seluruh platform.
// Operational Context adalah konteks bisnis yang WAJIB ada di setiap transaksi:
// - SchoolID: sekolah mana yang sedang aktif
// - AcademicYear: tahun ajaran (format "2026/2027")
// - Semester: semester (1 atau 2)
// - AcademicPeriodID: derivasi dari tahun+semester (format "20261")
// - TenantID: untuk multi-tenancy
// - UserID: user yang melakukan transaksi
// - Role: role user
package context

import (
	"fmt"
	"regexp"
)

// OperationalContext adalah konteks bisnis yang wajib di setiap transaksi.
// Struktur ini SYMMETRIC dengan OperationalContext di frontend JavaScript.
type OperationalContext struct {
	SchoolID         string `json:"schoolId"`
	SchoolName       string `json:"schoolName,omitempty"`
	AcademicYear     string `json:"academicYear"`
	Semester         int    `json:"semester"`
	AcademicPeriodID string `json:"academicPeriodId"`
	TenantID         string `json:"tenantId,omitempty"`
	UserID           string `json:"userId"`
	Role             string `json:"role"`
}

// academicYearPattern memvalidasi format "YYYY/YYYY".
var academicYearPattern = regexp.MustCompile(`^\d{4}/\d{4}$`)

// New membuat OperationalContext baru dengan validasi.
func New(schoolID, academicYear string, semester int, userID, role string) (*OperationalContext, error) {
	ctx := &OperationalContext{
		SchoolID:     schoolID,
		AcademicYear: academicYear,
		Semester:     semester,
		UserID:       userID,
		Role:         role,
	}

	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	ctx.AcademicPeriodID = ctx.DerivePeriodID()
	return ctx, nil
}

// DerivePeriodID menurunkan AcademicPeriodID dari AcademicYear + Semester.
// Contoh: "2026/2027" + 1 → "20261"
func (c *OperationalContext) DerivePeriodID() string {
	if len(c.AcademicYear) < 4 {
		return ""
	}
	return fmt.Sprintf("%s%d", c.AcademicYear[:4], c.Semester)
}

// Validate memvalidasi OperationalContext.
func (c *OperationalContext) Validate() error {
	if c.SchoolID == "" {
		return fmt.Errorf("schoolId wajib diisi")
	}
	if !academicYearPattern.MatchString(c.AcademicYear) {
		return fmt.Errorf("academicYear harus format YYYY/YYYY, got: %s", c.AcademicYear)
	}
	if c.Semester != 1 && c.Semester != 2 {
		return fmt.Errorf("semester harus 1 atau 2, got: %d", c.Semester)
	}
	if c.UserID == "" {
		return fmt.Errorf("userId wajib diisi")
	}
	if c.Role == "" {
		return fmt.Errorf("role wajib diisi")
	}
	return nil
}

// IsValid mengecek apakah context valid.
func (c *OperationalContext) IsValid() bool {
	return c.Validate() == nil
}

// DisplayLabel mengembalikan label untuk ditampilkan di UI.
func (c *OperationalContext) DisplayLabel() string {
	name := c.SchoolName
	if name == "" {
		name = c.SchoolID
	}
	return fmt.Sprintf("%s | %s - Semester %d", name, c.AcademicYear, c.Semester)
}

// Equals mengecek apakah dua context sama (berdasarkan school + period).
func (c *OperationalContext) Equals(other *OperationalContext) bool {
	if other == nil {
		return false
	}
	return c.SchoolID == other.SchoolID && c.AcademicPeriodID == other.AcademicPeriodID
}

// ToMap mengkonversi ke map (untuk logging, tracing, dll).
func (c *OperationalContext) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"schoolId":         c.SchoolID,
		"schoolName":       c.SchoolName,
		"academicYear":     c.AcademicYear,
		"semester":         c.Semester,
		"academicPeriodId": c.AcademicPeriodID,
		"tenantId":         c.TenantID,
		"userId":           c.UserID,
		"role":             c.Role,
	}
}
