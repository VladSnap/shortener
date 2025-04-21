package handlers

import (
	"context"
	"time"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	"go.uber.org/zap"
)

// Эти константы нужны, чтобы линтер не ругался на магические числа. В конфиг не вижу смысла выносить это.
const flushToDBIntervalSec = 5
const fanInChanSize = 10

// DeleterWorkerImpl - Реализация воркера для удаления сокращенных ссылок интерфейса DeleterWorker.
type DeleterWorkerImpl struct {
	fanInChan      chan chan services.DeleteShortID
	shorterService ShorterService
	buffer         []services.DeleteShortID
}

// NewDeleteWorker - Создает новую структуру DeleterWorkerImpl с указателем.
func NewDeleteWorker(shorterService ShorterService) *DeleterWorkerImpl {
	return &DeleterWorkerImpl{
		fanInChan:      make(chan chan services.DeleteShortID, fanInChanSize),
		shorterService: shorterService,
	}
}

// RunWork - Запускает горутину воркера, которая в фоне выполняет удаление сокращенных ссылок.
func (worker *DeleterWorkerImpl) RunWork() {
	// Таймер сброса буфера сообщений.
	ticker := time.NewTicker(flushToDBIntervalSec * time.Second)

	go func() {
		for {
			select {
			case deleteShorts := <-worker.fanInChan:
				for ds := range deleteShorts {
					worker.buffer = append(worker.buffer, ds)
				}
			case <-ticker.C:
				if len(worker.buffer) == 0 {
					continue
				}

				err := worker.shorterService.DeleteBatch(context.Background(), worker.buffer)
				if err != nil {
					log.Zap.Error("failed DeleteBatch", zap.Error(err))
					continue
				}

				worker.buffer = nil
			}
		}
	}()
}

// AddToDelete - Добавляет канал с идентификаторами сокращенных ссылок для
// потокобезопасного удаления используя паттерн FanIn.
func (worker *DeleterWorkerImpl) AddToDelete(shortIDs chan services.DeleteShortID) {
	worker.fanInChan <- shortIDs
}

// Close - Останавливает воркер, чтобы остановить приложение по graceful shutdown.
func (worker *DeleterWorkerImpl) Close() error {
	close(worker.fanInChan)
	return nil
}
