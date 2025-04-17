package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/VladSnap/shortener/internal/log"
	"go.uber.org/zap"
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
		bytes, err := w.zw.Write(b)
		if err != nil {
			return bytes, fmt.Errorf("failed gzip write: %w", err)
		}
		return bytes, nil
	}

	// Не сжимаем ответ
	bytes, err := w.ResponseWriter.Write(b)
	if err != nil {
		return bytes, fmt.Errorf("failed http write: %w", err)
	}
	return bytes, nil
}

func (w *gzipWriter) WriteHeader(statusCode int) {
	ct := w.Header().Get("Content-Type")

	// Надо проверить какой у нас контент будет в качестве ответа, чтобы принять решение, надо ли сжимать данные
	if ct != "" && slices.Contains(*gzipContentTypes, ct) && (statusCode < 300 || statusCode > 399) {
		w.Header().Add("Content-Encoding", "gzip")
		w.isCompressed = true
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipWriter) Close() error {
	if w.isCompressed {
		err := w.zw.Close()
		if err != nil {
			return fmt.Errorf("failed gzip close: %w", err)
		}
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
			defer func() {
				err := gzReader.Close()
				if err != nil {
					log.Zap.Error("failed gzip reader close", zap.Error(err))
				}
			}()
		}

		headerAE := r.Header.Get("Accept-Encoding")
		isAcceptGzip := strings.Contains(headerAE, "gzip")

		if isAcceptGzip {
			gzWriter := gzip.NewWriter(w)
			gzipWritterWrap := gzipWriter{w, gzWriter, false}
			ow = &gzipWritterWrap

			defer func() {
				err := gzipWritterWrap.Close()
				if err != nil {
					log.Zap.Error("failed gzip writer close", zap.Error(err))
				}
			}()
		}

		next.ServeHTTP(ow, r)
	})
}
