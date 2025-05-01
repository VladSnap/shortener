// Package auth реализует логику аутентификации пользователя с применением cookies.
package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// CookieAuthData - Структура для хранения информации о пользователе в cookie в формате JSON.
type CookieAuthData struct {
	UserID string `json:"user_id"`
}

const cookieValidSegmentCount int = 2

// CreateSignedCookie - Создает безопасную (подписаную) куки.
func CreateSignedCookie(userID, authKey string) (string, error) {
	data := CookieAuthData{userID}
	var jsonBuf bytes.Buffer
	err := json.NewEncoder(&jsonBuf).Encode(data)
	if err != nil {
		return "", fmt.Errorf("failed json encode: %w", err)
	}
	cookieContent := base64.RawURLEncoding.EncodeToString(jsonBuf.Bytes())
	// Подписываем алгоритмом HMAC, используя SHA-256.
	authCookieKey := sha256.Sum256([]byte(authKey))
	h := hmac.New(sha256.New, authCookieKey[:])
	h.Write([]byte(cookieContent))
	hmacCookie := h.Sum(nil)

	cookieSign := base64.RawURLEncoding.EncodeToString(hmacCookie)
	cookie := cookieContent + "." + cookieSign

	return cookie, nil
}

// VerifySignCookie - Проверяет подпись куки.
func VerifySignCookie(cookieValue string, authKey string) (bool, error) {
	cookieSegments := strings.Split(cookieValue, ".")
	if len(cookieSegments) != cookieValidSegmentCount {
		return false, errors.New("the cookie structure is not correct")
	}

	cookieContent := cookieSegments[0]
	cookieSign, err := base64.RawURLEncoding.DecodeString(cookieSegments[1])
	if err != nil {
		return false, fmt.Errorf("failed decode cookieSign from base64: %w", err)
	}

	authCookieKey := sha256.Sum256([]byte(authKey))
	h := hmac.New(sha256.New, authCookieKey[:])
	h.Write([]byte(cookieContent))
	hmacCookie := h.Sum(nil)

	isVerify := hmac.Equal(hmacCookie, cookieSign)
	if !isVerify {
		return isVerify, errors.New("hmac sign not equal")
	}
	return isVerify, nil
}

// DecodeCookie - Декодирует куки, чтобы далее извлечь из неё структуру CookieAuthData.
func DecodeCookie(cookieValue string) (*CookieAuthData, error) {
	cookieSegments := strings.Split(cookieValue, ".")
	if len(cookieSegments) != cookieValidSegmentCount {
		return nil, errors.New("the cookie structure is not correct")
	}

	cookieContent, err := base64.RawURLEncoding.DecodeString(cookieSegments[0])
	if err != nil {
		return nil, fmt.Errorf("failed decode cookieContent from base64: %w", err)
	}

	var cookieData CookieAuthData
	err = json.Unmarshal(cookieContent, &cookieData)
	if err != nil {
		return nil, fmt.Errorf("failed parsing cookie from json: %w", err)
	}

	return &cookieData, nil
}
