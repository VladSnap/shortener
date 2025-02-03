package repos

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/data/models"
	"github.com/VladSnap/shortener/internal/helpers"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/google/uuid"
)

type FileShortLinkRepo struct {
	links       map[string]models.ShortLinkData
	storageFile *os.File
}

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

func (repo *FileShortLinkRepo) CreateShortLink(shortID string, fullURL string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed create random: %w", err)
	}

	data := models.ShortLinkData{UUID: id.String(), ShortURL: shortID, OriginalURL: fullURL}
	repo.links[shortID] = data
	return repo.writeLink(data)
}

func (repo *FileShortLinkRepo) GetURL(shortID string) string {
	if data, ok := repo.links[shortID]; ok {
		return data.OriginalURL
	}
	return ""
}

func (repo *FileShortLinkRepo) Close() error {
	err := repo.storageFile.Close()
	if err != nil {
		return fmt.Errorf("file storage close error: %w", err)
	}
	log.Zap.Info("File storage closed")

	return nil
}

func (repo *FileShortLinkRepo) loadFromFile() ([]models.ShortLinkData, error) {
	scanner := bufio.NewScanner(repo.storageFile)
	var dataList []models.ShortLinkData

	for scanner.Scan() {
		dataBytes := scanner.Bytes()
		data := models.ShortLinkData{}
		err := json.Unmarshal(dataBytes, &data)
		if err != nil {
			return nil, fmt.Errorf("failed deserialize ShortLinkData: %w", err)
		}

		dataList = append(dataList, data)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed file scan: %w", err)
	}

	return dataList, nil
}

func (repo *FileShortLinkRepo) writeLink(link models.ShortLinkData) error {
	data, err := json.Marshal(link)
	if err != nil {
		return fmt.Errorf("failed serialize ShortLinkData: %w", err)
	}

	data = append(data, '\n')
	_, err = repo.storageFile.Write(data)
	if err != nil {
		return fmt.Errorf("failed write to file storage: %w", err)
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

func (repo *FileShortLinkRepo) loadLinks() (map[string]models.ShortLinkData, error) {
	links, err := repo.loadFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed load from file storage: %w", err)
	}
	linkMap := make(map[string]models.ShortLinkData, len(links))
	for _, link := range links {
		linkMap[link.ShortURL] = link
	}
	return linkMap, nil
}
