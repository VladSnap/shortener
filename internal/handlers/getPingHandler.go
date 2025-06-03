package handlers

import (
	"net/http"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	"go.uber.org/zap"
)

// GetPingHandler - Обработчик запроса проверки доступности базы данных.
type GetPingHandler struct {
	opts          *config.Options
	healthService *services.HealthService
}

// NewGetPingHandler - Создает новую структуру GetPingHandler с указателем.
func NewGetPingHandler(opts *config.Options) *GetPingHandler {
	handler := new(GetPingHandler)
	handler.opts = opts
	handler.healthService = services.NewHealthService(opts.DataBaseConnString)
	return handler
}

// Handle - Обрабатывает входящий запрос.
func (handler *GetPingHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, ValidateErrHTTPNotGET, http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	err := handler.healthService.PingDatabase(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte("OK"))
	if err != nil {
		log.Zap.Error(ErrFailedWriteToResponse, zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
