package handlers

import (
	"context"
	"time"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
)

const flushToDBIntervalSec = 5
const fanInChanSize = 10

type DeleterWorkerImpl struct {
	fanInChan      chan chan services.DeleteShortID
	shorterService ShorterService
	buffer         []services.DeleteShortID
}

func NewDeleteWorker(shorterService ShorterService) *DeleterWorkerImpl {
	return &DeleterWorkerImpl{
		fanInChan:      make(chan chan services.DeleteShortID, fanInChanSize),
		shorterService: shorterService,
	}
}

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
					log.Zap.Errorf("failed DeleteBatch: %w", err)
					continue
				}

				worker.buffer = nil
			}
		}
	}()
}

func (worker *DeleterWorkerImpl) AddToDelete(shortIDs chan services.DeleteShortID) {
	worker.fanInChan <- shortIDs
}

func (worker *DeleterWorkerImpl) Close() error {
	close(worker.fanInChan)
	return nil
}
