package validation

import (
	"errors"

	"github.com/VladSnap/shortener/internal/config"
)

type OptionsValidator struct {
}

func (vld *OptionsValidator) Validate(opts *config.Options) error {
	runesBaseURL := []rune(opts.BaseURL)

	if string(runesBaseURL[len(runesBaseURL)-1:]) == "/" {
		return errors.New("Incorrect -b argument. Don't put a slash at the end")
	}

	return nil
}
