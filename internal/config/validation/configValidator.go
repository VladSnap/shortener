package validation

import (
	"errors"

	"github.com/VladSnap/shortener/internal/config"
)

// OptionsValidator - Структура валидатора конфигов.
type OptionsValidator struct {
}

// Validate - Проверяет корректность конфига.
func (vld *OptionsValidator) Validate(opts *config.Options) error {
	runesBaseURL := []rune(opts.BaseURL)

	if string(runesBaseURL[len(runesBaseURL)-1:]) == "/" {
		return errors.New("incorrect -b argument, don't put a slash at the end")
	}

	return nil
}
