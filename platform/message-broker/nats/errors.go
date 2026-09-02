package nats

import "errors"

var (
	errEmptyURL         = errors.New("NATS URL tidak boleh kosong")
	errInvalidTimeout   = errors.New("timeout harus > 0")
	errNotConnected     = errors.New("tidak terhubung ke NATS server")
	errStreamNotFound   = errors.New("stream tidak ditemukan")
	errConsumerNotFound = errors.New("consumer tidak ditemukan")
	errPublishFailed    = errors.New("gagal publish message")
	errSubscribeFailed  = errors.New("gagal subscribe")
)

// ConnectionError adalah error spesifik connection.
type ConnectionError struct {
	URL   string
	Cause error
}

func (e *ConnectionError) Error() string {
	return "gagal terhubung ke NATS server " + e.URL + ": " + e.Cause.Error()
}

func (e *ConnectionError) Unwrap() error {
	return e.Cause
}
