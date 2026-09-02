// Package response menyediakan envelope standar untuk semua response API.
// Format ini SYMMETRIC dengan response dari Mock Adapter frontend.
//
// Format:
//
//	{
//	  "success": true/false,
//	  "data": {...} | null,
//	  "error": {code, message, details} | null
//	}
package response

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/platform/shared/errors"
)

// Response adalah envelope standar untuk semua response API.
// Struktur ini HARUS identik dengan response dari Mock Adapter frontend.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *ErrorInfo  `json:"error"`
}

// ErrorInfo berisi detail error.
// Struktur ini identik dengan error dari Mock Adapter frontend.
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success membuat response sukses dengan data.
func Success(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
		Error:   nil,
	}
}

// SuccessWithMeta membuat response sukses dengan data dan metadata (pagination, dll).
func SuccessWithMeta(data interface{}, meta interface{}) *Response {
	return &Response{
		Success: true,
		Data: map[string]interface{}{
			"items": data,
			"meta":  meta,
		},
		Error: nil,
	}
}

// Error membuat response error dari AppError.
func Error(err *errors.AppError) *Response {
	return &Response{
		Success: false,
		Data:    nil,
		Error: &ErrorInfo{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	}
}

// ErrorFromMessage membuat response error dari message sederhana.
// Gunakan hanya untuk error internal, bukan untuk error bisnis.
func ErrorFromMessage(code, message string) *Response {
	return &Response{
		Success: false,
		Data:    nil,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}

// Empty membuat response sukses tanpa data (untuk operasi delete, dll).
func Empty() *Response {
	return &Response{
		Success: true,
		Data:    nil,
		Error:   nil,
	}
}

// ToJSON mengkonversi response ke JSON bytes.
func (r *Response) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// ToJSONIndent mengkonversi response ke JSON bytes dengan indentasi.
func (r *Response) ToJSONIndent() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// HTTPStatus mengembalikan HTTP status code yang sesuai berdasarkan response.
func (r *Response) HTTPStatus() int {
	if r.Success {
		return http.StatusOK
	}

	if r.Error == nil {
		return http.StatusInternalServerError
	}

	// Map error code ke HTTP status
	switch r.Error.Code {
	case errors.CodeValidation, errors.CodeInvalidContext:
		return http.StatusBadRequest
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeDuplicate:
		return http.StatusConflict
	case errors.CodeUnauthorized:
		return http.StatusUnauthorized
	case errors.CodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
