package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// StatsResponse - Структура ответа для GetStatsHandler.
type StatsResponse struct {
	// CorrelationID - Количество сокращённых URL в сервисе.
	Urls int `json:"urls"`
	// ShortURL - Количество пользователей в сервисе.
	Users int `json:"users"`
}

// GetStatsHandler - .
type GetStatsHandler struct {
	opts    *config.Options
	service ShorterService
}

// NewGetStatsHandler - .
func NewGetStatsHandler(opts *config.Options, service ShorterService) *GetStatsHandler {
	handler := new(GetStatsHandler)
	handler.opts = opts
	handler.service = service
	return handler
}

// Handle - Обрабатывает входящий запрос.
func (handler *GetStatsHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, ValidateErrHTTPNotGET, http.StatusBadRequest)
		return
	}

	stats, err := handler.service.GetStats(req.Context())

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)

	result := StatsResponse{Urls: stats.Urls, Users: stats.Users}
	err = json.NewEncoder(res).Encode(result)

	if err != nil {
		log.Zap.Error(ErrFailedWriteToResponse, zap.Error(err))
		return
	}
}
