package middlewares

import (
	"net/http"
	"time"

	"github.com/VladSnap/shortener/internal/log"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Zap.Infof("Request %v %v", r.Method, r.RequestURI)
		// перед началом выполнения функции сохраняем текущее время
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		// вызываем следующий обработчик
		next.ServeHTTP(&lw, r)
		// после завершения замеряем время выполнения запроса
		duration := time.Since(start)

		log.Zap.Infoln(
			"Response", r.Method, r.RequestURI,
			"status:", responseData.status, // получаем перехваченный код статуса ответа
			"duration:", duration.Milliseconds(), "ms",
			"size:", responseData.size, "bytes", // получаем перехваченный размер ответа
		)
	})
}
