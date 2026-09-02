// Package client menyediakan gRPC client untuk komunikasi dengan tu-core.
package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	platformctx "sekolah-platform/platform/context"
)

// GRPCClient adalah client untuk tu-core gRPC services.
type GRPCClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewGRPCClient membuat gRPC client baru.
func NewGRPCClient(address string, timeout time.Duration) (*GRPCClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024), // 10MB
			grpc.MaxCallSendMsgSize(10*1024*1024),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("gagal connect ke tu-core: %w", err)
	}

	return &GRPCClient{
		conn:    conn,
		timeout: timeout,
	}, nil
}

// Conn mengembalikan underlying gRPC connection.
func (c *GRPCClient) Conn() *grpc.ClientConn {
	return c.conn
}

// Close menutup koneksi.
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

// ContextWithTimeout membuat context dengan timeout.
func (c *GRPCClient) ContextWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, c.timeout)
}

// InjectOperationalContext menginject OperationalContext ke gRPC metadata.
func InjectOperationalContext(ctx context.Context, opCtx *platformctx.OperationalContext) context.Context {
	md := metadata.Pairs(
		"x-school-id", opCtx.SchoolID,
		"x-school-name", opCtx.SchoolName,
		"x-academic-year", opCtx.AcademicYear,
		"x-semester", fmt.Sprintf("%d", opCtx.Semester),
		"x-academic-period", opCtx.AcademicPeriodID,
		"x-user-id", opCtx.UserID,
		"x-user-role", opCtx.Role,
	)
	return metadata.NewOutgoingContext(ctx, md)
}

// InjectCorrelation menginject correlation IDs ke gRPC metadata.
func InjectCorrelation(ctx context.Context, correlationID, traceID, requestID string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md = md.Copy()
	if correlationID != "" {
		md.Set("x-correlation-id", correlationID)
	}
	if traceID != "" {
		md.Set("x-trace-id", traceID)
	}
	if requestID != "" {
		md.Set("x-request-id", requestID)
	}
	return metadata.NewOutgoingContext(ctx, md)
}
