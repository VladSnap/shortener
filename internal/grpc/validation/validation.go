// Package validation содержит функции валидации для gRPC запросов.
package validation

import (
	"context"
	"fmt"
	"net/url"
	"unicode/utf8"

	"github.com/VladSnap/shortener/internal/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validationFailedErr = "validation failed: %w"

// ValidateOriginalURL проверяет корректность оригинальной URL.
func ValidateOriginalURL(originalURL string) error {
	if originalURL == "" {
		return fmt.Errorf(validationFailedErr, status.Error(codes.InvalidArgument, "original_url is required"))
	}

	// Парсим URL
	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return fmt.Errorf(validationFailedErr, status.Errorf(codes.InvalidArgument,
			"invalid original_url format: %v", err))
	}

	// Проверяем наличие схемы и хоста
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf(validationFailedErr, status.Error(codes.InvalidArgument,
			"original_url must contain schema and host"))
	}

	return nil
}

// ValidateShortID проверяет корректность короткого ID.
func ValidateShortID(shortID string) error {
	if shortID == "" {
		return fmt.Errorf(validationFailedErr, status.Error(codes.InvalidArgument, "short_id is required"))
	}

	if utf8.RuneCountInString(shortID) != constants.ShortIDLength {
		return fmt.Errorf(validationFailedErr, status.Errorf(codes.InvalidArgument,
			"short_id length must be %d", constants.ShortIDLength))
	}

	return nil
}

// ValidateShortURL проверяет корректность короткой URL.
func ValidateShortURL(shortURL string) error {
	if shortURL == "" {
		return fmt.Errorf(validationFailedErr, status.Error(codes.InvalidArgument, "short_url cannot be empty"))
	}

	if utf8.RuneCountInString(shortURL) != constants.ShortIDLength {
		return fmt.Errorf(validationFailedErr, status.Errorf(codes.InvalidArgument,
			"short_url length must be %d", constants.ShortIDLength))
	}

	return nil
}

// ValidateShortURLs проверяет массив коротких URL.
func ValidateShortURLs(shortURLs []string) error {
	if len(shortURLs) == 0 {
		return fmt.Errorf(validationFailedErr, status.Error(codes.InvalidArgument, "short_urls cannot be empty"))
	}

	for _, shortURL := range shortURLs {
		if err := ValidateShortURL(shortURL); err != nil {
			return fmt.Errorf(validationFailedErr, err)
		}
	}

	return nil
}

// ExtractUserID извлекает userID из контекста.
func ExtractUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return "", fmt.Errorf("context extraction failed: %w", status.Error(codes.Internal, "user ID not found in context"))
	}
	return userID, nil
}

// ValidateContextDeadline проверяет, не истек ли контекст.
func ValidateContextDeadline(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context validation failed: %w", status.Error(codes.DeadlineExceeded, "context deadline exceeded"))
	default:
		return nil
	}
}
