package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/VladSnap/shortener/internal/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// HealthService provides health check functionality for the application.
type HealthService struct {
	dbConnString string
}

// NewHealthService creates a new HealthService instance.
func NewHealthService(dbConnString string) *HealthService {
	return &HealthService{
		dbConnString: dbConnString,
	}
}

// PingDatabase checks if the database connection is healthy.
func (s *HealthService) PingDatabase(ctx context.Context) error {
	if s.dbConnString == "" {
		return nil // No database configured, consider healthy
	}

	db, err := sql.Open("pgx", s.dbConnString)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Zap.Error("failed to close database connection", zap.Error(closeErr))
		}
	}()

	const timeOutPingSec = 5
	pingCtx, cancel := context.WithTimeout(ctx, timeOutPingSec*time.Second)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
