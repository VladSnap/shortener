// Package validation хранит логику валидации для обработчиков http запросов.
package validation

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/VladSnap/shortener/internal/constants"
)

// ValidateShortURL - Валидирует сокращенную ссылку.
func ValidateShortURL(inputURL string) error {
	if inputURL == "" {
		return errors.New("shortURL should not be empty")
	}
	if utf8.RuneCountInString(inputURL) != constants.ShortIDLength {
		return errors.New("shortURL length should be 8")
	}
	return nil
}

// ValidateURL - Валидирует оригинальную ссылку.
func ValidateURL(inputURL string, paramName string) error {
	// Проверяем, что строка не пустая
	if inputURL == "" {
		return fmt.Errorf("%s can't be empty", paramName)
	}
	// Парсим URL
	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return fmt.Errorf("incorrect format %s: %w", paramName, err)
	}
	// Проверяем наличие схемы и хоста
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("%s must contain schema and host", paramName)
	}

	return nil
}

// ValidatePath - Валидирует path ссылки.
func ValidatePath(path string) bool {
	segments := strings.Split(path, "/")
	return len(segments) == 2 && segments[1] != ""
}
