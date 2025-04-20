package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/VladSnap/shortener/internal/handlers/mocks"
	"github.com/VladSnap/shortener/internal/services"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
				responseBody: "Http method not GET\n",
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
				responseBody: "Request path incorrect\n",
				location:     "",
			},
		},
		{
			name:        "request path not corrected #2",
			requestPath: "/foo/bar",
			httpMethod:  http.MethodGet,
			id:          "dVCBBnmd",
			url:         "",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Request path incorrect\n",
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
				responseBody: "Request path incorrect\n",
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
				code:         404,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Url not found\n",
				location:     "",
			},
		},
	}

	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := m.NewMockShorterService(ctrl)
	getHandler := NewGetHandler(mockService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slink *services.ShortedLink
			if tt.url == "" {
				slink = nil
			} else {
				slink = &services.ShortedLink{OriginalURL: tt.url}
			}
			mockService.EXPECT().GetURL(ctx, tt.id).
				Return(slink, nil).
				AnyTimes()

			request := httptest.NewRequest(tt.httpMethod, tt.requestPath, http.NoBody)
			request.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()
			getHandler.Handle(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "no error for read response")
			err = res.Body.Close()
			assert.NoError(t, err, "no error for close response body")

			assert.Equal(t, tt.want.contentType, res.Header.Get(HeaderContentType))
			assert.Equal(t, tt.url, res.Header.Get("Location"))
			assert.NotEmpty(t, string(resBody))
		})
	}
}
