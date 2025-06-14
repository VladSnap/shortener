package interceptors

import (
	"testing"
)

func TestShouldValidateMethod(t *testing.T) {
	tests := []struct {
		name         string
		fullMethod   string
		config       TrustedSubnetConfig
		expectResult bool
	}{
		{
			name:       "empty protected methods - validate all",
			fullMethod: "/shortener.ShortenerService/CreateShortLink",
			config: TrustedSubnetConfig{
				ProtectedMethods: []string{},
			},
			expectResult: true,
		},
		{
			name:       "exact method match",
			fullMethod: "/shortener.ShortenerService/GetStats",
			config: TrustedSubnetConfig{
				ProtectedMethods: []string{"/shortener.ShortenerService/GetStats"},
				UseMethodSuffix:  false,
			},
			expectResult: true,
		},
		{
			name:       "method name without service prefix",
			fullMethod: "/shortener.ShortenerService/GetStats",
			config: TrustedSubnetConfig{
				ProtectedMethods: []string{"GetStats"},
				UseMethodSuffix:  false,
			},
			expectResult: true,
		},
		{
			name:       "suffix match",
			fullMethod: "/shortener.ShortenerService/GetStats",
			config: TrustedSubnetConfig{
				ProtectedMethods: []string{"Stats"},
				UseMethodSuffix:  true,
			},
			expectResult: true,
		},
		{
			name:       "no match - different method",
			fullMethod: "/shortener.ShortenerService/CreateShortLink",
			config: TrustedSubnetConfig{
				ProtectedMethods: []string{"GetStats"},
				UseMethodSuffix:  false,
			},
			expectResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldValidateMethod(tt.fullMethod, tt.config)
			if result != tt.expectResult {
				t.Errorf("shouldValidateMethod() = %v, want %v", result, tt.expectResult)
			}
		})
	}
}
