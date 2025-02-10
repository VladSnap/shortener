package data

import "fmt"

type DuplicateShortLinkError struct {
	ShortURL string
}

func NewDuplicateError(shortURL string) error {
	return &DuplicateShortLinkError{
		ShortURL: shortURL,
	}
}

func (de *DuplicateShortLinkError) Error() string {
	return fmt.Sprintf("shortURL '%v' already exists in storage", de.ShortURL)
}
