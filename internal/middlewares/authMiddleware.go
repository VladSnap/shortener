package middlewares

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type cookieAuthData struct {
	UserID string `json:"user_id"`
}

// Такое следует хранить в защищенном хранилище, а не в коде или в конфигах.
var authCookieKey = sha256.Sum256([]byte("kl1jmo;u6hn&*On0jo8f"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("Auth")
		if err != nil {
			userID, err := setNewAuthCookie(w)
			if err != nil {
				log.Zap.Warnf("failed setNewAuthCookie: %w", err)
			}
			r = r.WithContext(context.WithValue(r.Context(), constants.UserIDContextKey, userID))
		} else {
			if _, err := verifySignCookie(authCookie.Value); err != nil {
				log.Zap.Warnf("failed verifySignCookie: %w", err)
				_, err = setNewAuthCookie(w)
				if err != nil {
					log.Zap.Warnf("failed setNewAuthCookie: %w", err)
				}
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authData, err := decodeCookie(authCookie.Value)
			if err != nil {
				fmt.Printf("err %w", err)
				http.Error(w, "Not decoded cookie: %w", http.StatusInternalServerError)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), constants.UserIDContextKey, authData.UserID))
		}

		next.ServeHTTP(w, r)
	})
}

func setNewAuthCookie(w http.ResponseWriter) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed generate new user id: %w", err)
	}

	cookieValue, err := createSignedCookie(id.String())
	if err != nil {
		return "", fmt.Errorf("failed createSignedCookie: %w", err)
	}

	cookie := &http.Cookie{Name: "Auth", Value: cookieValue, Path: "/"}
	http.SetCookie(w, cookie)
	return id.String(), nil
}

func createSignedCookie(userID string) (string, error) {
	data := cookieAuthData{userID}
	var jsonBuf bytes.Buffer
	err := json.NewEncoder(&jsonBuf).Encode(data)
	if err != nil {
		return "", fmt.Errorf("failed json encode: %w", err)
	}
	cookieContent := base64.RawURLEncoding.EncodeToString(jsonBuf.Bytes())
	// Подписываем алгоритмом HMAC, используя SHA-256.
	h := hmac.New(sha256.New, authCookieKey[:])
	h.Write([]byte(cookieContent))
	hmacCookie := h.Sum(nil)

	cookieSign := base64.RawURLEncoding.EncodeToString(hmacCookie)
	cookie := cookieContent + "." + cookieSign

	return cookie, nil
}

func verifySignCookie(cookieValue string) (bool, error) {
	cookieSegments := strings.Split(cookieValue, ".")
	if len(cookieSegments) != 2 {
		return false, errors.New("the cookie structure is not correct")
	}

	cookieContent := cookieSegments[0]
	cookieSign, err := base64.RawURLEncoding.DecodeString(cookieSegments[1])
	if err != nil {
		return false, fmt.Errorf("failed decode cookieSign from base64", err)
	}

	h := hmac.New(sha256.New, authCookieKey[:])
	h.Write([]byte(cookieContent))
	hmacCookie := h.Sum(nil)

	iserify := hmac.Equal(hmacCookie, cookieSign)
	if !iserify {
		return iserify, errors.New("hmac sign not equal")
	}
	return iserify, nil
}

func decodeCookie(cookieValue string) (*cookieAuthData, error) {
	cookieSegments := strings.Split(cookieValue, ".")
	if len(cookieSegments) != 2 {
		return nil, errors.New("the cookie structure is not correct")
	}

	cookieContent, err := base64.RawURLEncoding.DecodeString(cookieSegments[0])
	if err != nil {
		return nil, fmt.Errorf("failed decode cookieContent from base64", err)
	}

	var cookieData cookieAuthData
	err = json.Unmarshal(cookieContent, &cookieData)
	if err != nil {
		return nil, fmt.Errorf("failed parsing cookie from json", err)
	}

	return &cookieData, nil
}
