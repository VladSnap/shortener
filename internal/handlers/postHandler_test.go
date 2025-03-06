package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VladSnap/shortener/internal/constants"
	m "github.com/VladSnap/shortener/internal/handlers/mocks"
	"github.com/VladSnap/shortener/internal/services"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL string = "http://localhost:8080"

func TestPostHandler(t *testing.T) {
	type want struct {
		code         int
		contentType  string
		responseBody string
	}
	tests := []struct {
		name        string
		httpMethod  string
		requestPath string
		sourceURL   string
		shortID     string
		want        want
	}{
		{
			name:        "positive test #1",
			sourceURL:   "http://test.url",
			requestPath: "/",
			httpMethod:  http.MethodPost,
			shortID:     "fVdpTFBo",
			want: want{
				code:         201,
				contentType:  "text/plain",
				responseBody: baseURL + "/fVdpTFBo",
			},
		}, {
			name:        "request url invalid",
			httpMethod:  http.MethodPost,
			requestPath: "/",
			sourceURL:   "google.com",
			shortID:     "sKbYvAgT",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "req.Body verify error\n",
			},
		}, {
			name:        "request body is empty",
			httpMethod:  http.MethodPost,
			requestPath: "/",
			sourceURL:   "",
			shortID:     "sKbYvAgT",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Required url\n",
			},
		}, {
			name:        "http method not corrected",
			httpMethod:  http.MethodGet,
			requestPath: "/",
			sourceURL:   "",
			shortID:     "sVpHyErn",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Http method not POST\n",
			},
		}, {
			name:        "request path not correct #1",
			sourceURL:   "http://test3.url",
			requestPath: "//",
			httpMethod:  http.MethodPost,
			shortID:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Incorrect request path\n",
			},
		}, {
			name:        "request path not correct #2",
			sourceURL:   "http://test4.url",
			requestPath: "/foo",
			httpMethod:  http.MethodPost,
			shortID:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Incorrect request path\n",
			},
		}, {
			name:        "request path not correct #3",
			sourceURL:   "http://test5.url",
			requestPath: "/foo/bar",
			httpMethod:  http.MethodPost,
			shortID:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Incorrect request path\n",
			},
		}, {
			name:        "internal server error",
			sourceURL:   "http://test6.url",
			requestPath: "/",
			httpMethod:  http.MethodPost,
			shortID:     "fbUhNtPv",
			want: want{
				code:         500,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "random fail\n",
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := m.NewMockShorterService(ctrl)
	postHandler := NewPostHandler(mockService, baseURL)
	ret := &services.ShortedLink{URL: ""}
	userID := "d1a8485a-430a-49f4-92ba-50886e1b07c6"
	ctx := context.WithValue(context.Background(), constants.UserIDContextKey, userID)
	mockService.EXPECT().CreateShortLink(ctx, "http://test6.url", userID).
		Return(ret, errors.New("random fail")).
		AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret = &services.ShortedLink{URL: tt.shortID}
			mockService.EXPECT().CreateShortLink(ctx, tt.sourceURL, userID).
				Return(ret, nil).
				AnyTimes()

			r := strings.NewReader(tt.sourceURL)
			postRequest := httptest.NewRequest(tt.httpMethod, tt.requestPath, r)
			postRequest.Header.Add(HeaderContentType, "text/plain; charset=utf-8")
			postRequest = postRequest.WithContext(ctx)
			w := httptest.NewRecorder()
			postHandler.Handle(w, postRequest)
			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode, "Incorrect status code")
			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "no error for read response")
			err = res.Body.Close()
			assert.NoError(t, err, "no error for close response body")

			shortURL := string(resBody)
			assert.Equal(t, tt.want.contentType, res.Header.Get(HeaderContentType), "Incorrect header content-type")
			assert.Equal(t, tt.want.responseBody, shortURL, "Incorrect response short url")
		})
	}
}
