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
	zw           *gzip.Writer
	isCompressed bool
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	if w.isCompressed {
		// Сжимаем ответ, если у него подходящий тип контента
		return w.zw.Write(b)
	}

	// Не сжимаем ответ
	return w.ResponseWriter.Write(b)
}

func (w *gzipWriter) WriteHeader(statusCode int) {
	ct := w.Header().Get("Content-Type")

	// Надо проверить какой у нас контент будет в качестве ответа, чтобы принять решение, надо ли сжимать данные
	if ct != "" && slices.Contains(*gzipContentTypes, ct) && statusCode < 300 {
		w.Header().Add("Content-Encoding", "gzip")
		w.isCompressed = true
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipWriter) Close() error {
	if w.isCompressed {
		return w.zw.Close()
	}

	return nil
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
			gzipWritterWrap := gzipWriter{w, gzWriter, false}
			ow = &gzipWritterWrap
			defer gzipWritterWrap.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
