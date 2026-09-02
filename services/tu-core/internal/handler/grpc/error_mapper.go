package grpc

import (
	"sekolah-platform/platform/shared/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapAppErrorToGRPCStatus mengkonversi AppError ke gRPC status.
func MapAppErrorToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	appErr, ok := err.(*errors.AppError)
	if !ok {
		// Bukan AppError, return sebagai internal error
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}

	var code codes.Code
	switch appErr.Code {
	case errors.CodeValidation, errors.CodeInvalidContext:
		code = codes.InvalidArgument
	case errors.CodeNotFound:
		code = codes.NotFound
	case errors.CodeDuplicate:
		code = codes.AlreadyExists
	case errors.CodeUnauthorized, errors.CodeTokenExpired:
		code = codes.Unauthenticated
	case errors.CodeForbidden:
		code = codes.PermissionDenied
	case errors.CodeBusinessRule:
		code = codes.FailedPrecondition
	case errors.CodeConflict:
		code = codes.Aborted
	case errors.CodeTimeout:
		code = codes.DeadlineExceeded
	case errors.CodeIntegration:
		code = codes.Unavailable
	default:
		code = codes.Internal
	}

	return status.Errorf(code, "[%s] %s", appErr.Code, appErr.Message)
}
