package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"

	//"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	type want struct {
		code         int
		contentType  string
		responseBody string
		location     string
	}
	tests := []struct {
		name        string
		httpMethod  string
		requestPath string
		id          string
		url         string
		want        want
	}{
		{
			name:        "positive test #1",
			requestPath: "/{id}",
			httpMethod:  http.MethodGet,
			id:          "fVjYdBgR",
			url:         "http://test.url",
			want: want{
				code:         307,
				contentType:  "text/html; charset=utf-8",
				responseBody: "<a href=\"http://test.url\">Temporary Redirect</a>.\n\n",
				location:     "http://test.url",
			},
		},
		{
			name:        "positive test #2",
			requestPath: "/{id}",
			httpMethod:  http.MethodGet,
			id:          "cDkNYhTB",
			url:         "http://test2.url",
			want: want{
				code:         307,
				contentType:  "text/html; charset=utf-8",
				responseBody: "<a href=\"http://test2.url\">Temporary Redirect</a>.\n\n",
				location:     "http://test2.url",
			},
		},
		{
			name:        "http method not corrected",
			requestPath: "/{id}",
			httpMethod:  http.MethodPost,
			id:          "PdjMBGtd",
			url:         "",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
				location:     "",
			},
		},
		{
			name:        "request path not corrected #1",
			requestPath: "/{id}/foo",
			httpMethod:  http.MethodGet,
			id:          "sonYHbTD",
			url:         "",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
				location:     "",
			},
		},
		{
			name:        "request path not corrected #2",
			requestPath: "/foo",
			httpMethod:  http.MethodGet,
			id:          "dVCBBnmd",
			url:         "",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
				location:     "",
			},
		},
		{
			name:        "request path value empty",
			requestPath: "/{id}",
			httpMethod:  http.MethodGet,
			id:          "",
			url:         "",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
				location:     "",
			},
		},
		{
			name:        "url not found",
			requestPath: "/{id}",
			httpMethod:  http.MethodGet,
			id:          "bdGTBvoP",
			url:         "",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
				location:     "",
			},
		},
	}

	shortLinkRepo := new(MockShortLinkRepo)
	getHandler := NewGetHandler(shortLinkRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortLinkRepo.On("GetURL", tt.id).Return(tt.url)

			request := httptest.NewRequest(tt.httpMethod, tt.requestPath, nil)
			request.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()
			getHandler.Handle(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.url, res.Header.Get("Location"))
			assert.NotEmpty(t, string(resBody))
		})
	}
}
