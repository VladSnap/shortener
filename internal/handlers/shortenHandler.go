package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/validation"
	"go.uber.org/zap"
)

// ShortenRequest - Структура запроса для ShortenHandler.
type ShortenRequest struct {
	// URL - Оригинальный URL который требуется сократить.
	URL string `json:"url"`
}

// ShortenResponse - Структура ответа для ShortenHandler.
type ShortenResponse struct {
	// Result - Результат в виде сокращенной ссылки.
	Result string `json:"result"`
}

// ShortenHandler - Обработчик запроса сокращения одной ссылки в формате json.
type ShortenHandler struct {
	service ShorterService
	baseURL string
}

// NewShortenHandler - Создает новую структуру ShortenHandler с указателем.
func NewShortenHandler(service ShorterService, baseURL string) *ShortenHandler {
	handler := new(ShortenHandler)
	handler.service = service
	handler.baseURL = baseURL
	return handler
}

// Handle - Обрабатывает входящий запрос.
func (handler *ShortenHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Http method not POST", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get(HeaderContentType)

	if !strings.Contains(ct, HeaderApplicationJSONValue) && !strings.Contains(ct, HeaderApplicationXgzipValue) {
		http.Error(res, "Incorrect content-type:"+ct, http.StatusBadRequest)
		return
	}

	var request ShortenRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if request.URL == "" {
		http.Error(res, "Required url", http.StatusBadRequest)
		return
	}

	request.URL = strings.TrimSuffix(request.URL, "\r")
	request.URL = strings.TrimSuffix(request.URL, "\n")
	if err := validation.ValidateURL(request.URL, "URL"); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	userID := ""
	if value, ok := req.Context().Value(constants.UserIDContextKey).(string); ok {
		userID = value
	}
	shortLink, err := handler.service.CreateShortLink(req.Context(), request.URL, userID)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result := ShortenResponse{Result: handler.baseURL + "/" + shortLink.URL}

	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)
	if shortLink.IsDuplicated {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	err = json.NewEncoder(res).Encode(result)

	if err != nil {
		log.Zap.Error(ErrFailedWriteToResponse, zap.Error(err))
		return
	}
}
