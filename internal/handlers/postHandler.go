package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/validation"
	"go.uber.org/zap"
)

// PostHandler - Обработчик запроса сокращения одной ссылки в формате text/plain.
type PostHandler struct {
	service ShorterService
	baseURL string
}

// NewPostHandler - Создает новую структуру PostHandler с указателем.
func NewPostHandler(service ShorterService, baseURL string) *PostHandler {
	handler := new(PostHandler)
	handler.service = service
	handler.baseURL = baseURL
	return handler
}

// Handle - Обрабатывает входящий запрос.
func (handler *PostHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Http method not POST", http.StatusBadRequest)
		return
	}

	if req.URL.Path != "/" {
		http.Error(res, "Incorrect request path", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get(HeaderContentType)
	if !strings.Contains(ct, "text/plain") && !strings.Contains(ct, HeaderApplicationXgzipValue) {
		http.Error(res, "Incorrect content-type:"+ct, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	fullURL := string(body)

	if err != nil || fullURL == "" {
		http.Error(res, "Required url", http.StatusBadRequest)
		return
	}
	fullURL = strings.TrimSuffix(fullURL, "\r")
	fullURL = strings.TrimSuffix(fullURL, "\n")
	if err := validation.ValidateURL(fullURL, "req.Body"); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	userID := ""
	if value, ok := req.Context().Value(constants.UserIDContextKey).(string); ok {
		userID = value
	}
	shortLink, err := handler.service.CreateShortLink(req.Context(), fullURL, userID)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Add(HeaderContentType, "text/plain")
	if shortLink.IsDuplicated {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	_, err = res.Write([]byte(handler.baseURL + "/" + shortLink.URL))

	if err != nil {
		log.Zap.Error(ErrFailedWriteToResponse, zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
