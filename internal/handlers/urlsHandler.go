package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"go.uber.org/zap"
)

// ShortedLinkResponse - Структура ответа для UrlsHandler.
type ShortedLinkResponse struct {
	// OriginalURL - Оригинальный URL который был сокращен.
	OriginalURL string `json:"original_url"`
	// ShortURL - Сокращенная ссылка.
	ShortURL string `json:"short_url"`
}

// UrlsHandler - Обработчик запроса чтения сокращенных ссылок пользователя.
type UrlsHandler struct {
	service ShorterService
	baseURL string
}

// NewUrlsHandler - Создает новую структуру UrlsHandler с указателем.
func NewUrlsHandler(service ShorterService, baseURL string) *UrlsHandler {
	handler := new(UrlsHandler)
	handler.service = service
	handler.baseURL = baseURL
	return handler
}

// Handle - Обрабатывает входящий запрос.
func (handler *UrlsHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Http method not GET", http.StatusBadRequest)
		return
	}

	userID := ""
	if value, ok := req.Context().Value(constants.UserIDContextKey).(string); ok {
		userID = value
	}
	shortedLinks, err := handler.service.GetAllByUserID(req.Context(), userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if len(shortedLinks) == 0 {
		http.Error(res, "Urls for user not found", http.StatusNoContent)
		return
	}

	responseRows := make([]*ShortedLinkResponse, 0, len(shortedLinks))
	for _, sl := range shortedLinks {
		rr := &ShortedLinkResponse{sl.OriginalURL, handler.baseURL + "/" + sl.URL}
		responseRows = append(responseRows, rr)
	}

	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(responseRows)

	if err != nil {
		log.Zap.Error(ErrFailedWriteToResponse, zap.Error(err))
		return
	}
}
