package repos

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/helpers"
	"github.com/VladSnap/shortener/internal/log"
	"golang.org/x/exp/maps"
)

// FileShortLinkRepo - Репозиторий для доступа к файловому хранилищу сокращателя ссылок.
type FileShortLinkRepo struct {
	links       map[string]*data.ShortLinkData
	storageFile *os.File
}

// NewFileShortLinkRepo - Создает новую структуру FileShortLinkRepo с указателем.
func NewFileShortLinkRepo(fileStoragePath string) (*FileShortLinkRepo, error) {
	repo := new(FileShortLinkRepo)
	file, err := createFileStorage(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed create file storage: %w", err)
	}
	repo.storageFile = file

	links, err := repo.loadLinks()
	if err != nil {
		return nil, fmt.Errorf("failed load links: %w", err)
	}
	repo.links = links

	return repo, nil
}

// Add - Сохраняет структуру сокращенной ссылки в файле.
func (repo *FileShortLinkRepo) Add(ctx context.Context, link *data.ShortLinkData) (
	*data.ShortLinkData, error) {
	repo.links[link.ShortURL] = link
	err := repo.writeLink(link)
	if err != nil {
		return nil, fmt.Errorf("failed write link to file storage: %w", err)
	}

	return link, nil
}

// AddBatch - Сохраняет пачку структур сокращенных ссылок в файле.
func (repo *FileShortLinkRepo) AddBatch(ctx context.Context, links []*data.ShortLinkData) (
	[]*data.ShortLinkData, error) {
	for _, link := range links {
		repo.links[link.ShortURL] = link
	}

	err := repo.writeLinkBatch(links)
	if err != nil {
		return nil, fmt.Errorf("failed write batch links to file storage: %w", err)
	}

	return links, nil
}

// Get - Читает полную ссылку по сокращенной ссылке.
func (repo *FileShortLinkRepo) Get(ctx context.Context, shortID string) (*data.ShortLinkData, error) {
	link := repo.links[shortID]
	return link, nil
}

// GetAllByUserID - Получить все сокращенные ссылки указанного пользователя.
func (repo *FileShortLinkRepo) GetAllByUserID(ctx context.Context, userID string) (
	[]*data.ShortLinkData, error) {
	var links []*data.ShortLinkData

	for _, l := range repo.links {
		if l.UserID == userID {
			links = append(links, l)
		}
	}

	return links, nil
}

// DeleteBatch - Удаляет пачку структур сокращенных ссылок в БД.
func (repo *FileShortLinkRepo) DeleteBatch(ctx context.Context, shortIDs []data.DeleteShortData) error {
	// Сначала обновляем записи в мемори кэше.
	for _, sid := range shortIDs {
		link := repo.links[sid.ShortURL]
		if link.UserID == sid.UserID {
			link.IsDeleted = true
		}
	}
	// Удаляем содержимое файла для перезаписи.
	err := repo.storageFile.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed truncate file storage: %w", err)
	}
	// Перезаписываем содержимое файла, чтобы проставить флаг is_deleted.
	_, err = repo.AddBatch(ctx, maps.Values(repo.links))
	if err != nil {
		return fmt.Errorf("failed rewrite file after batch delete: %w", err)
	}
	return nil
}

// GetStats - Получает статистику о пользователях и всех ссылках.
func (repo *FileShortLinkRepo) GetStats(ctx context.Context) (*data.StatsData, error) {
	data := data.NewStatsData(len(repo.links), repo.calcAllUsers())
	return data, nil
}

// Close - Закрывает файл.
func (repo *FileShortLinkRepo) Close() error {
	err := repo.storageFile.Close()
	if err != nil {
		return fmt.Errorf("file storage close error: %w", err)
	}
	log.Zap.Info("File storage closed")

	return nil
}

func (repo *FileShortLinkRepo) loadFromFile() ([]*data.ShortLinkData, error) {
	scanner := bufio.NewScanner(repo.storageFile)
	var dataList []*data.ShortLinkData

	for scanner.Scan() {
		dataBytes := scanner.Bytes()
		sd := data.ShortLinkData{}
		err := json.Unmarshal(dataBytes, &sd)
		if err != nil {
			return nil, fmt.Errorf("failed deserialize ShortLinkData: %w", err)
		}

		dataList = append(dataList, &sd)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed file scan: %w", err)
	}

	return dataList, nil
}

func (repo *FileShortLinkRepo) writeLink(link *data.ShortLinkData) error {
	writer := bufio.NewWriter(repo.storageFile)
	sd, err := json.Marshal(link)
	if err != nil {
		return fmt.Errorf("failed serialize ShortLinkData: %w", err)
	}

	if _, err := writer.Write(sd); err != nil {
		return fmt.Errorf("failed write to file buffer: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("failed write \\n to file buffer: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed flush buffer to fileStorage: %w", err)
	}
	return nil
}

func (repo *FileShortLinkRepo) writeLinkBatch(links []*data.ShortLinkData) error {
	writer := bufio.NewWriter(repo.storageFile)
	for _, link := range links {
		sd, err := json.Marshal(link)
		if err != nil {
			return fmt.Errorf("failed serialize ShortLinkData: %w", err)
		}

		if _, err := writer.Write(sd); err != nil {
			return fmt.Errorf("failed write to file buffer: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("failed write \\n to file buffer: %w", err)
		}
	}

	err := writer.Flush()
	if err != nil {
		return fmt.Errorf("failed flush buffer to fileStorage: %w", err)
	}
	return nil
}

func createFileStorage(fileStoragePath string) (*os.File, error) {
	ok, err := helpers.DirectoryExists(filepath.Dir(fileStoragePath))
	if !ok && err == nil {
		return nil, errors.New("fileStoragePath directory not exists")
	} else if err != nil {
		return nil, fmt.Errorf("failed check fileStoragePath: %w", err)
	}

	file, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, constants.FileRWPerm)
	if err != nil {
		return nil, fmt.Errorf("failed open file storage: %w", err)
	}

	return file, nil
}

func (repo *FileShortLinkRepo) loadLinks() (map[string]*data.ShortLinkData, error) {
	links, err := repo.loadFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed load from file storage: %w", err)
	}
	linkMap := make(map[string]*data.ShortLinkData, len(links))
	for _, link := range links {
		linkMap[link.ShortURL] = link
	}
	return linkMap, nil
}

func (repo *FileShortLinkRepo) calcAllUsers() int {
	users := make(map[string]bool)

	for _, l := range repo.links {
		users[l.UserID] = true
	}

	return len(users)
}
