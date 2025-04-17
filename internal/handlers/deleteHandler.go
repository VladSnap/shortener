package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/services"
	"github.com/VladSnap/shortener/internal/validation"
)

type DeleteHandler struct {
	deleteWorker DeleterWorker
}

//go:generate mockgen -destination=mocks/deleteWorker_mock.go -package=mocks github.com/VladSnap/shortener/internal/handlers DeleterWorker

type DeleterWorker interface {
	Close() error
	AddToDelete(shortIDs chan services.DeleteShortID)
	RunWork()
}

func NewDeleteHandler(deleteWorker DeleterWorker) *DeleteHandler {
	handler := new(DeleteHandler)
	handler.deleteWorker = deleteWorker
	return handler
}

func (handler *DeleteHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Http method not DELETE", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get(HeaderContentType)
	if !strings.Contains(ct, HeaderApplicationJSONValue) && !strings.Contains(ct, HeaderApplicationXgzipValue) {
		http.Error(res, "Incorrect content-type:"+ct, http.StatusBadRequest)
		return
	}

	var shortURLs []string
	if err := json.NewDecoder(req.Body).Decode(&shortURLs); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	for _, surl := range shortURLs {
		if err := validation.ValidateShortURL(surl); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if len(shortURLs) == 0 {
		res.WriteHeader(http.StatusNotAcceptable)
		return
	}

	userID := ""
	if value, ok := req.Context().Value(constants.UserIDContextKey).(string); ok {
		userID = value
	}

	const toDeleteChanSize = 100

	toDeleteChan := make(chan services.DeleteShortID, toDeleteChanSize)
	defer close(toDeleteChan)

	const cancelTimeoutSec = 2
	// Создаем контекст с таймаутом 2 секунды
	ctx, cancel := context.WithTimeout(context.Background(), cancelTimeoutSec*time.Second)
	defer cancel()

rng:
	for _, url := range shortURLs {
		select {
		case <-ctx.Done():
			break rng // Выйдем из цикла, если мы не уложились в таймаут записи данных, канал автоматически закроется.
		default:
			deleteSID := services.NewDeleteShortID(url, userID)
			toDeleteChan <- deleteSID
		}
	}

	handler.deleteWorker.AddToDelete(toDeleteChan)

	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)
	res.WriteHeader(http.StatusAccepted)
}
