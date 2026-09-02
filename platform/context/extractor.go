package context

import (
	"fmt"
	"strconv"
)

// Header names untuk HTTP request.
const (
	HeaderSchoolID = "X-School-Id"
	HeaderSchoolName = "X-School-Name"
	HeaderAcademicYear = "X-Academic-Year"
	HeaderSemester = "X-Semester"
	HeaderAcademicPeriod = "X-Academic-Period"
	HeaderTenantID = "X-Tenant-Id"
	HeaderUserID = "X-User-Id"
	HeaderUserRole = "X-User-Role"
)

// HeaderProvider adalah interface untuk penyedia header (HTTP, gRPC, dll).
type HeaderProvider interface {
	Get(key string) string
}

// MapProvider adalah HeaderProvider dari map (untuk testing).
type MapProvider map[string]string

// Get mengimplementasikan HeaderProvider.
func (m MapProvider) Get(key string) string {
	return m[key]
}

// Extract mengekstrak OperationalContext dari header provider.
func Extract(provider HeaderProvider) (*OperationalContext, error) {
	schoolID := provider.Get(HeaderSchoolID)
	schoolName := provider.Get(HeaderSchoolName)
	academicYear := provider.Get(HeaderAcademicYear)
	semesterStr := provider.Get(HeaderSemester)
	periodID := provider.Get(HeaderAcademicPeriod)
	tenantID := provider.Get(HeaderTenantID)
	userID := provider.Get(HeaderUserID)
	role := provider.Get(HeaderUserRole)

	// Jika periodID sudah ada, parse academicYear dan semester darinya
	if periodID != "" && academicYear == "" {
		academicYear, semester := ParsePeriodID(periodID)
		if academicYear != "" {
			return &OperationalContext{
				SchoolID:         schoolID,
				SchoolName:       schoolName,
				AcademicYear:     academicYear,
				Semester:         semester,
				AcademicPeriodID: periodID,
				TenantID:         tenantID,
				UserID:           userID,
				Role:             role,
			}, nil
		}
	}

	// Parse semester
	semester := 0
	if semesterStr != "" {
		var err error
		semester, err = strconv.Atoi(semesterStr)
		if err != nil {
			return nil, fmt.Errorf("semester tidak valid: %s", semesterStr)
		}
	}

	ctx := &OperationalContext{
		SchoolID:     schoolID,
		SchoolName:   schoolName,
		AcademicYear: academicYear,
		Semester:     semester,
		TenantID:     tenantID,
		UserID:       userID,
		Role:         role,
	}

	// Derive periodID jika belum ada
	if ctx.AcademicPeriodID == "" && ctx.AcademicYear != "" && ctx.Semester > 0 {
		ctx.AcademicPeriodID = ctx.DerivePeriodID()
	}

	return ctx, nil
}

// ParsePeriodID mengurai periodID menjadi academicYear dan semester.
// Contoh: "20261" → ("2026/2027", 1)
func ParsePeriodID(periodID string) (string, int) {
	if len(periodID) < 5 {
		return "", 0
	}
	yearStr := periodID[:4]
	semesterStr := periodID[4:]
	semester, err := strconv.Atoi(semesterStr)
	if err != nil {
		return "", 0
	}
	// Derive end year
	year := 0
	for _, c := range yearStr {
		year = year*10 + int(c-'0')
	}
	endYear := year + 1
	return fmt.Sprintf("%s/%d", yearStr, endYear), semester
}

// ToHeaders mengkonversi OperationalContext ke map headers.
func (c *OperationalContext) ToHeaders() map[string]string {
	headers := map[string]string{
		HeaderSchoolID:       c.SchoolID,
		HeaderAcademicYear:   c.AcademicYear,
		HeaderSemester:       strconv.Itoa(c.Semester),
		HeaderAcademicPeriod: c.AcademicPeriodID,
		HeaderUserID:         c.UserID,
		HeaderUserRole:       c.Role,
	}
	if c.SchoolName != "" {
		headers[HeaderSchoolName] = c.SchoolName
	}
	if c.TenantID != "" {
		headers[HeaderTenantID] = c.TenantID
	}
	return headers
}