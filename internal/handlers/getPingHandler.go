package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/log"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type GetPingHandler struct {
	opts *config.Options
}

func NewGetPingHandler(opts *config.Options) *GetPingHandler {
	handler := new(GetPingHandler)
	handler.opts = opts
	return handler
}

func (handler *GetPingHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Http method not GET", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("pgx", handler.opts.DataBaseConnString)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Zap.Error("failed database connection close: %w", err)
		}
	}()

	const timeOutPingSec = 5
	ctx, cancel := context.WithTimeout(context.Background(), timeOutPingSec*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte("OK"))
	if err != nil {
		log.Zap.Errorf("failed write to response: %w", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
