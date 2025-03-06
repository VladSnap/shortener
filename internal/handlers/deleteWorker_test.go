package handlers_test

import (
	"errors"
	"testing"
	"time"

	"github.com/VladSnap/shortener/internal/handlers"
	m "github.com/VladSnap/shortener/internal/handlers/mocks"
	"github.com/VladSnap/shortener/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const flushToDBIntervalSec = 6

func TestDeleterWorker_AddToDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShorterService := m.NewMockShorterService(ctrl)
	worker := handlers.NewDeleteWorker(mockShorterService)

	shortIDs := make(chan services.DeleteShortID, 5)
	shortIDs <- services.DeleteShortID{ShortURL: "url1", UserID: "user1"}
	close(shortIDs)

	// Мокируем DeleteBatch для проверки передачи данных
	mockShorterService.EXPECT().
		DeleteBatch(gomock.Any(), []services.DeleteShortID{
			{ShortURL: "url1", UserID: "user1"},
		}).
		Return(nil).
		Times(1)

	worker.AddToDelete(shortIDs)

	// Запускаем RunWork
	go worker.RunWork()

	// Ждем срабатывания таймера (используем фиксированное время вместо приватной константы)
	time.Sleep(flushToDBIntervalSec * time.Second) // Примерно больше, чем FlushToDBIntervalSec
}

func TestDeleterWorker_RunWork_FlushToDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShorterService := m.NewMockShorterService(ctrl)
	worker := handlers.NewDeleteWorker(mockShorterService)

	// Мокируем DeleteBatch
	mockShorterService.EXPECT().
		DeleteBatch(gomock.Any(), []services.DeleteShortID{
			{ShortURL: "url1", UserID: "user1"},
		}).
		Return(nil).
		Times(1)

	shortIDs := make(chan services.DeleteShortID, 5)
	shortIDs <- services.DeleteShortID{ShortURL: "url1", UserID: "user1"}
	close(shortIDs)

	worker.AddToDelete(shortIDs)

	// Запускаем RunWork
	go worker.RunWork()

	// Ждем срабатывания таймера (используем фиксированное время вместо приватной константы)
	time.Sleep(flushToDBIntervalSec * time.Second) // Примерно больше, чем FlushToDBIntervalSec
}

func TestDeleterWorker_RunWork_DeleteBatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShorterService := m.NewMockShorterService(ctrl)
	worker := handlers.NewDeleteWorker(mockShorterService)

	// Первый вызов DeleteBatch завершается ошибкой
	mockShorterService.EXPECT().
		DeleteBatch(gomock.Any(), []services.DeleteShortID{
			{ShortURL: "url1", UserID: "user1"},
		}).
		Return(errors.New("database error")).
		Times(1)

	// Второй вызов DeleteBatch успешен (повторная попытка после ошибки)
	mockShorterService.EXPECT().
		DeleteBatch(gomock.Any(), []services.DeleteShortID{
			{ShortURL: "url1", UserID: "user1"},
		}).
		Return(nil).
		Times(1)

	shortIDs := make(chan services.DeleteShortID, 5)
	shortIDs <- services.DeleteShortID{ShortURL: "url1", UserID: "user1"}
	close(shortIDs)

	worker.AddToDelete(shortIDs)

	// Запускаем RunWork
	go worker.RunWork()

	// Ждем срабатывания таймера дважды (используем фиксированное время вместо приватной константы)
	time.Sleep(flushToDBIntervalSec * 2 * time.Second) // Два интервала по ~6 секунд
}

func TestDeleterWorker_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShorterService := m.NewMockShorterService(ctrl)
	worker := handlers.NewDeleteWorker(mockShorterService)

	err := worker.Close()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем, что канал закрыт, отправляя данные в него
	require.Panics(t, func() { worker.AddToDelete(make(chan services.DeleteShortID)) }, "expected no panic")
}
