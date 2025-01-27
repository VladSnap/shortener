package middlewares

import (
	"log"
	"net/http"
	"time"
)

func TimerTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// перед началом выполнения функции сохраняем текущее время
		start := time.Now()
		// вызываем следующий обработчик
		next.ServeHTTP(w, r)
		// после завершения замеряем время выполнения запроса
		duration := time.Since(start)

		log.Printf("Request %v %v handled in %s", r.Method, r.RequestURI, duration)
	})
}
