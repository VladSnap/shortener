package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VladSnap/shortener/internal/log"
	"go.uber.org/zap"
)

type (
	// Берём структуру для хранения сведений об ответе.
	responseData struct {
		data   string
		status int
		size   int
	}

	// Добавляем реализацию http.ResponseWriter.
	loggingResponseWriter struct {
		http.ResponseWriter // Встраиваем оригинальный http.ResponseWriter.
		responseData        *responseData
	}
)

// Write - Реализует метод записи тела ответа для обертки логирования.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// Записываем ответ, используя оригинальный http.ResponseWriter.
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // Захватываем размер.
	r.responseData.data += string(b)
	if err != nil {
		return size, fmt.Errorf("failed logging response write: %w", err)
	}
	return size, nil
}

// WriteHeader - Реализует метод записи заголовка для обертки логирования.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// Записываем код статуса, используя оригинальный http.ResponseWriter.
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // Захватываем код статуса.
}

// LogMiddleware - Мидлварь для логирования запросов и ответов.
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rqHeaders := ""
		for k, v := range r.Header {
			rqHeaders += fmt.Sprintf("%s: %v | ", k, v)
		}
		// Перед началом выполнения функции сохраняем текущее время.
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // Встраиваем оригинальный http.ResponseWriter.
			responseData:   responseData,
		}
		var duration time.Duration
		rsHeaders := ""

		defer func() {
			duration = time.Since(start)
			for k, v := range w.Header() {
				rsHeaders += fmt.Sprintf("%s: %v | ", k, v)
			}
			log.Zap.Info("Request",
				zap.String("path", r.Method+" "+r.RequestURI),
				zap.Int("status", responseData.status),
				zap.Duration("duration", duration),
				zap.Int("byte_size", responseData.size),
				zap.String("request_headers", rqHeaders),
				zap.String("response_headers", rsHeaders),
				zap.String("data", responseData.data),
			)
		}()

		// Вызываем следующий обработчик.
		next.ServeHTTP(&lw, r)
		// После завершения замеряем время выполнения запроса.
	})
}
