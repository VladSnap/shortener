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
		return fmt.Errorf("required %s", paramName)
	}
	if !strings.Contains(inputURL, "http") {
		return fmt.Errorf("%s verify error", paramName)
	}
	if !strings.Contains(inputURL, "://") {
		return fmt.Errorf("%s verify error", paramName)
	}
	// Парсим URL
	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil || parsedURL.Host == "" {
		return fmt.Errorf("%s verify error", paramName)
	}

	return nil
}

// ValidatePath - Валидирует path ссылки.
func ValidatePath(path string) bool {
	segments := strings.Split(path, "/")
	return len(segments) == 2 && segments[1] != ""
}
