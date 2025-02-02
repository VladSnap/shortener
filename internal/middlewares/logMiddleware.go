package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VladSnap/shortener/internal/log"
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

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// Записываем ответ, используя оригинальный http.ResponseWriter.
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // Захватываем размер.
	r.responseData.data += string(b)
	return size, fmt.Errorf("failed logging response write: %w", err)
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// Записываем код статуса, используя оригинальный http.ResponseWriter.
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // Захватываем код статуса.
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Zap.Infof("Request %v %v", r.Method, r.RequestURI)
		headersLog := ""
		for k, v := range r.Header {
			headersLog += fmt.Sprintf("%s: %v | ", k, v)
		}
		log.Zap.Infof("Headers: %s", headersLog)
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

		// Вызываем следующий обработчик.
		next.ServeHTTP(&lw, r)
		// После завершения замеряем время выполнения запроса.
		duration := time.Since(start)

		headersLog = ""
		for k, v := range w.Header() {
			headersLog += fmt.Sprintf("%s: %v | ", k, v)
		}

		log.Zap.Infoln(
			"Response", r.Method, r.RequestURI,
			"status:", responseData.status, // Получаем перехваченный код статуса ответа.
			"duration:", duration.Milliseconds(), "ms",
			"size:", responseData.size, "bytes", // Получаем перехваченный размер ответа.
			"\nHeaders:", headersLog,
			"\nData:", "'"+responseData.data+"'",
		)
	})
}
