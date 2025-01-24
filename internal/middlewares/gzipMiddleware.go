package middlewares

import (
	"compress/gzip"
	"net/http"
	"slices"
	"strings"
)

var gzipContentTypes *[]string = &[]string{"application/json", "text/html"}

type gzipWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	ct := w.Header().Get("Content-Type")

	if ct != "" && slices.Contains(*gzipContentTypes, ct) {
		// Сжимаем ответ, если у него подходящий тип контента
		return w.zw.Write(b)
	}

	// Не сжимаем ответ
	return w.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		headerCE := r.Header.Get("Content-Encoding")
		isRequestGzip := strings.Contains(headerCE, "gzip")

		if isRequestGzip {
			gzReader, err := gzip.NewReader(r.Body)

			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			r.Body = gzReader
			defer gzReader.Close()
		}

		headerAE := r.Header.Get("Accept-Encoding")
		isAcceptGzip := strings.Contains(headerAE, "gzip")

		if isAcceptGzip {
			gzWriter := gzip.NewWriter(w)
			ow = &gzipWriter{w, gzWriter}
			defer gzWriter.Close()

			w.Header().Set("Content-Encoding", "gzip")
		}

		next.ServeHTTP(ow, r)
	})
}
