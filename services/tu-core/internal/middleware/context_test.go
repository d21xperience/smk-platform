package middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestOperationalContextUnaryInterceptor(t *testing.T) {
	interceptor := OperationalContextUnaryInterceptor()

	md := metadata.Pairs(
		MetaSchoolID, "smk-001",
		MetaAcademicYear, "2026/2027",
		MetaSemester, "1",
		MetaUserID, "user-123",
		MetaUserRole, "operator",
	)

	ctx := metadata.NewIncomingContext(context.Background(), md)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		opCtx := GetOperationalContext(ctx)
		if opCtx == nil {
			t.Error("Expected OperationalContext in context")
			return nil, nil
		}
		if opCtx.SchoolID != "smk-001" {
			t.Errorf("Expected schoolID 'smk-001', got '%s'", opCtx.SchoolID)
		}
		if opCtx.AcademicPeriodID != "20261" {
			t.Errorf("Expected periodID '20261', got '%s'", opCtx.AcademicPeriodID)
		}
		if opCtx.UserID != "user-123" {
			t.Errorf("Expected userID 'user-123', got '%s'", opCtx.UserID)
		}
		return "ok", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test"}
	resp, err := interceptor(ctx, nil, info, handler)

	if err != nil {
		t.Fatalf("Interceptor failed: %v", err)
	}
	if resp != "ok" {
		t.Error("Expected response 'ok'")
	}
}

func TestOperationalContextMissingMetadata(t *testing.T) {
	interceptor := OperationalContextUnaryInterceptor()

	// Context tanpa metadata
	ctx := context.Background()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test"}
	_, err := interceptor(ctx, nil, info, handler)

	if err == nil {
		t.Error("Expected error for missing metadata")
	}
}

func TestOperationalContextInvalidContext(t *testing.T) {
	interceptor := OperationalContextUnaryInterceptor()

	// Metadata tidak lengkap (missing schoolId)
	md := metadata.Pairs(
		MetaAcademicYear, "2026/2027",
		MetaSemester, "1",
	)

	ctx := metadata.NewIncomingContext(context.Background(), md)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test"}
	_, err := interceptor(ctx, nil, info, handler)

	if err == nil {
		t.Error("Expected error for invalid context")
	}
}

func TestOperationalContextFromPeriodID(t *testing.T) {
	interceptor := OperationalContextUnaryInterceptor()

	// Metadata dengan periodID saja
	md := metadata.Pairs(
		MetaSchoolID, "smk-001",
		MetaAcademicPeriod, "20261",
		MetaUserID, "user-123",
		MetaUserRole, "operator",
	)

	ctx := metadata.NewIncomingContext(context.Background(), md)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		opCtx := GetOperationalContext(ctx)
		if opCtx.AcademicYear != "2026/2027" {
			t.Errorf("Expected academicYear '2026/2027', got '%s'", opCtx.AcademicYear)
		}
		if opCtx.Semester != 1 {
			t.Errorf("Expected semester 1, got %d", opCtx.Semester)
		}
		return "ok", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test"}
	_, err := interceptor(ctx, nil, info, handler)

	if err != nil {
		t.Fatalf("Interceptor failed: %v", err)
	}
}
