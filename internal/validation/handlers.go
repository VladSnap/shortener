package validation

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/VladSnap/shortener/internal/constants"
	urlverifier "github.com/davidmytton/url-verifier"
)

func ValidateShortURL(url string) error {
	if url == "" {
		return errors.New("shortURL should not be empty")
	}
	if utf8.RuneCountInString(url) != constants.ShortIDLength {
		return errors.New("shortURL length should be 8")
	}
	return nil
}

func ValidateURL(url string, paramName string) error {
	if url == "" {
		return fmt.Errorf("required %s", paramName)
	}
	verifyRes, urlIsValid := urlverifier.NewVerifier().Verify(url)
	if urlIsValid != nil || !verifyRes.IsURL || !verifyRes.IsRFC3986URL {
		return fmt.Errorf("%s verify error", paramName)
	}
	return nil
}

func ValidatePath(path string) bool {
	segments := strings.Split(path, "/")
	return len(segments) == 2 && segments[1] != ""
}
