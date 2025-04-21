package data

import "fmt"

// DuplicateShortLinkError - Структура ошибки дублирования сокращенной ссылки.
type DuplicateShortLinkError struct {
	ShortURL string
}

// NewDuplicateError - Создает новую структуру DuplicateShortLinkError с указателем.
func NewDuplicateError(shortURL string) error {
	return &DuplicateShortLinkError{
		ShortURL: shortURL,
	}
}

// Error - Реализует интерфейс Error.
func (de *DuplicateShortLinkError) Error() string {
	return fmt.Sprintf("shortURL '%v' already exists in storage", de.ShortURL)
}
