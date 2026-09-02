package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Logger interface untuk logging.
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, args ...interface{})  { log.Printf("[INFO] "+msg, args...) }
func (l *defaultLogger) Error(msg string, args ...interface{}) { log.Printf("[ERROR] "+msg, args...) }
func (l *defaultLogger) Debug(msg string, args ...interface{}) { log.Printf("[DEBUG] "+msg, args...) }

// LoggingUnaryInterceptor log setiap gRPC request/response.
func LoggingUnaryInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = &defaultLogger{}
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Extract metadata untuk logging
		opCtx := GetOperationalContext(ctx)
		schoolID := ""
		userID := ""
		if opCtx != nil {
			schoolID = opCtx.SchoolID
			userID = opCtx.UserID
		}

		logger.Info("gRPC request started",
			"method", info.FullMethod,
			"school_id", schoolID,
			"user_id", userID,
		)

		// Call handler
		resp, err := handler(ctx, req)

		// Log result
		duration := time.Since(start)
		statusCode := status.Code(err)

		if err != nil {
			logger.Error("gRPC request failed",
				"method", info.FullMethod,
				"school_id", schoolID,
				"user_id", userID,
				"duration", duration,
				"status", statusCode.String(),
				"error", err.Error(),
			)
		} else {
			logger.Info("gRPC request completed",
				"method", info.FullMethod,
				"school_id", schoolID,
				"user_id", userID,
				"duration", duration,
				"status", statusCode.String(),
			)
		}

		return resp, err
	}
}
