package validation

import (
	"fmt"
	"strings"

	urlverifier "github.com/davidmytton/url-verifier"
)

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
