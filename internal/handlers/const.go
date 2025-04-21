package handlers

// Константы http заголовков и текста ошибок.
const (
	// ErrFailedWriteToResponse - Ошибка записи ответа.
	ErrFailedWriteToResponse = "failed write to response: %w"
	// HeaderContentType - Http заголовок Content-Type.
	HeaderContentType = "Content-Type"
	// HeaderApplicationJSONValue - Http заголовок application/json.
	HeaderApplicationJSONValue = "application/json"
	// HeaderApplicationXgzipValue - Http заголовок application/x-gzip.
	HeaderApplicationXgzipValue = "application/x-gzip"
)
