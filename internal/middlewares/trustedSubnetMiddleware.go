package middlewares

import (
	"net"
	"net/http"
)

// TrustedSubnetMiddleware - Мидлварь для авторизации по доверенной подсети.
func TrustedSubnetMiddleware(trustedSubnet string) func(next http.Handler) http.Handler {
	_, subnet, _ := net.ParseCIDR(trustedSubnet)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if subnet == nil {
				http.Error(w, "Trusted subnet not set", http.StatusForbidden)
				return
			}

			realIP := r.Header.Get("X-Real-IP")
			if realIP == "" {
				http.Error(w, "Not set X-Real-IP", http.StatusForbidden)
				return
			}

			ip := net.ParseIP(realIP)
			if ip == nil {
				http.Error(w, "Invalid IP", http.StatusForbidden)
				return
			}

			if subnet.Contains(ip) {
				// Вызываем следующий обработчик.
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "IP subnet not trusted", http.StatusForbidden)
				return
			}
		})
	}
}
