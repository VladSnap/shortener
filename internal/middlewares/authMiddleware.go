package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VladSnap/shortener/internal/auth"
	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthMiddleware - Мидлварь для аутентификации и атворизации пользователя.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("Auth")
		if err != nil {
			userID := handleMissingAuthCookie(w)
			r = r.WithContext(context.WithValue(r.Context(), constants.UserIDContextKey, userID))
			next.ServeHTTP(w, r)
			return
		}

		if !handleUnauthorized(w, authCookie) {
			return
		}

		authData, ok := handleAuthCookie(w, authCookie)

		if !ok {
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), constants.UserIDContextKey, authData.UserID))
		next.ServeHTTP(w, r)
	})
}

func handleMissingAuthCookie(w http.ResponseWriter) string {
	userID, err := setNewAuthCookie(w)
	if err != nil {
		log.Zap.Warn("failed setNewAuthCookie", zap.Error(err))
	}

	return userID
}

func handleUnauthorized(w http.ResponseWriter, authCookie *http.Cookie) bool {
	if _, err := auth.VerifySignCookie(authCookie.Value); err != nil {
		log.Zap.Warn("failed verifySignCookie", zap.Error(err))
		_, err = setNewAuthCookie(w)
		if err != nil {
			log.Zap.Warn("failed setNewAuthCookie", zap.Error(err))
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func handleAuthCookie(w http.ResponseWriter, authCookie *http.Cookie) (*auth.CookieAuthData, bool) {
	authData, err := auth.DecodeCookie(authCookie.Value)
	if err != nil {
		log.Zap.Warn("failed decodeCookie", zap.Error(err))
		http.Error(w, "Not decoded cookie: %w", http.StatusInternalServerError)
		return nil, false
	}
	return authData, true
}

func setNewAuthCookie(w http.ResponseWriter) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed generate new user id: %w", err)
	}

	cookieValue, err := auth.CreateSignedCookie(id.String())
	if err != nil {
		return "", fmt.Errorf("failed createSignedCookie: %w", err)
	}

	cookie := &http.Cookie{Name: "Auth", Value: cookieValue, Path: "/"}
	http.SetCookie(w, cookie)
	return id.String(), nil
}
