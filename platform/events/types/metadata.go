package types

// EventMetadata berisi metadata observability untuk event.
type EventMetadata struct {
	CorrelationID string `json:"correlationId"`
	TraceID       string `json:"traceId"`
	RequestID     string `json:"requestId"`
}

// EmptyMetadata membuat metadata kosong.
func EmptyMetadata() EventMetadata {
	return EventMetadata{}
}

// NewMetadata membuat metadata baru.
func NewMetadata(correlationID, traceID, requestID string) EventMetadata {
	return EventMetadata{
		CorrelationID: correlationID,
		TraceID:       traceID,
		RequestID:     requestID,
	}
}