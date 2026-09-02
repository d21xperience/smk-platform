// Package middleware menyediakan gRPC interceptors untuk cross-cutting concerns.
package middleware

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	platformctx "sekolah-platform/platform/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Metadata keys untuk Operational Context.
const (
	MetaSchoolID       = "x-school-id"
	MetaSchoolName     = "x-school-name"
	MetaAcademicYear   = "x-academic-year"
	MetaSemester       = "x-semester"
	MetaAcademicPeriod = "x-academic-period"
	MetaTenantID       = "x-tenant-id"
	MetaUserID         = "x-user-id"
	MetaUserRole       = "x-user-role"
	MetaCorrelationID  = "x-correlation-id"
	MetaTraceID        = "x-trace-id"
	MetaRequestID      = "x-request-id"
)

// contextKey adalah type untuk context keys.
type contextKey string

const operationalContextKey contextKey = "operational_context"

// Regex untuk validasi format academic year (YYYY/YYYY)
var academicYearRegex = regexp.MustCompile(`^\d{4}/\d{4}$`)

// OperationalContextUnaryInterceptor mengekstrak OperationalContext dari gRPC metadata
// dan inject ke context.Context.
func OperationalContextUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		// Build OperationalContext
		opCtx, err := extractOperationalContext(md)
		if err != nil {
			log.Printf("[Middleware] ❌ Invalid operational context: %v", err)
			return nil, status.Errorf(codes.InvalidArgument, "invalid operational context: %v", err)
		}

		log.Printf("[Middleware] ✅ Operational context loaded: school=%s, year=%s, semester=%d, user=%s",
			opCtx.SchoolID, opCtx.AcademicYear, opCtx.Semester, opCtx.UserID)

		// Inject ke context
		ctx = context.WithValue(ctx, operationalContextKey, opCtx)

		// Continue to handler
		return handler(ctx, req)
	}
}

// extractOperationalContext membangun OperationalContext dari gRPC metadata.
// PATCH: Validasi dilakukan secara lokal untuk bypass bug di platform/context.Validate()
func extractOperationalContext(md metadata.MD) (*platformctx.OperationalContext, error) {
	schoolID := getFirst(md, MetaSchoolID)
	schoolName := getFirst(md, MetaSchoolName)
	academicYear := getFirst(md, MetaAcademicYear)
	semesterStr := getFirst(md, MetaSemester)
	periodID := getFirst(md, MetaAcademicPeriod)
	tenantID := getFirst(md, MetaTenantID)
	userID := getFirst(md, MetaUserID)
	role := getFirst(md, MetaUserRole)

	// DEBUG: Log semua metadata yang diterima
	log.Printf("[Middleware] 📥 Raw metadata: schoolID=%q, academicYear=%q, semester=%q, periodID=%q, userID=%q, role=%q",
		schoolID, academicYear, semesterStr, periodID, userID, role)

	// Parse semester
	semester := 0
	if semesterStr != "" {
		var err error
		semester, err = strconv.Atoi(semesterStr)
		if err != nil {
			return nil, fmt.Errorf("invalid semester format: %w", err)
		}
	}

	// Jika periodID ada tapi academicYear tidak, derive dari periodID
	if periodID != "" && academicYear == "" {
		academicYear, semester = platformctx.ParsePeriodID(periodID)
	}

	ctx := &platformctx.OperationalContext{
		SchoolID:         schoolID,
		SchoolName:       schoolName,
		AcademicYear:     academicYear,
		Semester:         semester,
		AcademicPeriodID: periodID,
		TenantID:         tenantID,
		UserID:           userID,
		Role:             role,
	}

	// Derive periodID jika belum ada
	if ctx.AcademicPeriodID == "" && ctx.AcademicYear != "" && ctx.Semester > 0 {
		ctx.AcademicPeriodID = ctx.DerivePeriodID()
	}

	// ============================================================
	// PATCH: Validasi LOKAL (bypass platform/context.Validate())
	// ============================================================
	// Alasan: platform/context.Validate() memiliki bug yang menyebabkan
	// field AcademicYear terbaca kosong meskipun sudah di-set dengan benar.
	// Validasi lokal ini memberikan kontrol penuh dan debugging yang lebih baik.
	// ============================================================
	if err := validateOperationalContext(ctx); err != nil {
		return nil, err
	}

	return ctx, nil
}

// validateOperationalContext melakukan validasi lokal terhadap OperationalContext.
// Ini adalah pengganti untuk platformctx.OperationalContext.Validate() yang bermasalah.
func validateOperationalContext(ctx *platformctx.OperationalContext) error {
	// 1. SchoolID wajib
	if ctx.SchoolID == "" {
		return fmt.Errorf("schoolID wajib diisi")
	}

	// 2. AcademicYear wajib (jika tidak ada, derive dari periodID)
	if ctx.AcademicYear == "" {
		return fmt.Errorf("academicYear wajib diisi (format YYYY/YYYY, contoh: 2024/2025)")
	}

	// 3. Validasi format AcademicYear (YYYY/YYYY)
	if !academicYearRegex.MatchString(ctx.AcademicYear) {
		return fmt.Errorf("academicYear harus format YYYY/YYYY, got: %q", ctx.AcademicYear)
	}

	// 4. Semester wajib 1 atau 2
	if ctx.Semester != 1 && ctx.Semester != 2 {
		return fmt.Errorf("semester harus 1 atau 2, got: %d", ctx.Semester)
	}

	// 5. UserID wajib
	if ctx.UserID == "" {
		return fmt.Errorf("userID wajib diisi")
	}

	// 6. Role wajib
	if ctx.Role == "" {
		return fmt.Errorf("role wajib diisi")
	}

	// 7. AcademicPeriodID (optional, tapi jika ada harus valid format)
	if ctx.AcademicPeriodID != "" {
		// Format: YYYY-YYYY-S (contoh: 2024-2025-1)
		periodRegex := regexp.MustCompile(`^\d{4}-\d{4}-\d$`)
		if !periodRegex.MatchString(ctx.AcademicPeriodID) {
			return fmt.Errorf("academicPeriodID harus format YYYY-YYYY-S, got: %q", ctx.AcademicPeriodID)
		}
	}

	return nil
}

// GetOperationalContext mengambil OperationalContext dari context.Context.
func GetOperationalContext(ctx context.Context) *platformctx.OperationalContext {
	if v := ctx.Value(operationalContextKey); v != nil {
		if opCtx, ok := v.(*platformctx.OperationalContext); ok {
			return opCtx
		}
	}
	return nil
}

// MustGetOperationalContext mengambil OperationalContext atau panic.
func MustGetOperationalContext(ctx context.Context) *platformctx.OperationalContext {
	opCtx := GetOperationalContext(ctx)
	if opCtx == nil {
		panic("OperationalContext not found in context")
	}
	return opCtx
}

// getFirst mengambil value pertama dari metadata.
func getFirst(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// InjectOperationalContextToMetadata menginject OperationalContext ke gRPC metadata.
// Digunakan oleh client (edge-bff) saat memanggil tu-core.
func InjectOperationalContextToMetadata(ctx context.Context, opCtx *platformctx.OperationalContext) context.Context {
	md := metadata.Pairs(
		MetaSchoolID, opCtx.SchoolID,
		MetaSchoolName, opCtx.SchoolName,
		MetaAcademicYear, opCtx.AcademicYear,
		MetaSemester, strconv.Itoa(opCtx.Semester),
		MetaAcademicPeriod, opCtx.AcademicPeriodID,
		MetaTenantID, opCtx.TenantID,
		MetaUserID, opCtx.UserID,
		MetaUserRole, opCtx.Role,
	)
	return metadata.NewOutgoingContext(ctx, md)
}

// WithOperationalContext menyuntikkan OperationalContext ke dalam context.Context
// agar bisa dibaca oleh Service dan Domain Engine melalui GetOperationalContext.
func WithOperationalContext(ctx context.Context, opCtx *platformctx.OperationalContext) context.Context {
	return context.WithValue(ctx, operationalContextKey, opCtx)
}
