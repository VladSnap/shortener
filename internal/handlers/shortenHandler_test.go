package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/services"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenHandler(t *testing.T) {
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
			requestPath: "/api/shorten",
			httpMethod:  http.MethodPost,
			shortID:     "fVdpTFBo",
			want: want{
				code:         201,
				contentType:  HeaderApplicationJSONValue,
				responseBody: fmt.Sprintf("{\"result\":\"%v/fVdpTFBo\"}\n", baseURL),
			},
		}, {
			name:        "request url invalid",
			httpMethod:  http.MethodPost,
			requestPath: "/api/shorten",
			sourceURL:   "google.com",
			shortID:     "sKbYvAgT",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "URL verify error\n",
			},
		}, {
			name:        "request body is empty",
			httpMethod:  http.MethodPost,
			requestPath: "/api/shorten",
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
			requestPath: "/api/shorten",
			sourceURL:   "",
			shortID:     "sVpHyErn",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Http method not POST\n",
			},
		}, {
			name:        "internal server error",
			sourceURL:   "http://test6.url",
			requestPath: "/api/shorten",
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
	mockService := NewMockShorterService(ctrl)
	handler := NewShortenHandler(mockService, baseURL)
	ret := &services.ShortedLink{URL: ""}
	userID := "d1a8485a-430a-49f4-92ba-50886e1b07c6"
	ctx := context.WithValue(context.Background(), constants.UserIDContextKey, userID)
	mockService.EXPECT().CreateShortLink(ctx, "http://test6.url", userID).
		Return(ret, errors.New("random fail")).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret = &services.ShortedLink{URL: tt.shortID}
			mockService.EXPECT().CreateShortLink(ctx, tt.sourceURL, userID).
				Return(ret, nil).AnyTimes()

			requestData := ShortenRequest{
				URL: tt.sourceURL,
			}
			rqBytes, err := json.Marshal(requestData)
			if err != nil {
				assert.Error(t, err)
			}

			r := bytes.NewReader(rqBytes)
			postRequest := httptest.NewRequest(tt.httpMethod, tt.requestPath, r)
			postRequest.Header.Add(HeaderContentType, HeaderApplicationJSONValue)
			postRequest = postRequest.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.Handle(w, postRequest)
			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode, "Incorrect status code")

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "no error for read response")
			err = res.Body.Close()
			assert.NoError(t, err, "no error for close response body")
			resBodyStr := string(resBody)

			assert.Equal(t, tt.want.contentType, res.Header.Get(HeaderContentType), "Incorrect header content-type")
			assert.Equal(t, tt.want.responseBody, resBodyStr, "Incorrect response body")
		})
	}
}
