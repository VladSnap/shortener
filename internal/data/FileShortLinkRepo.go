package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/google/uuid"
)

type FileShortLinkRepo struct {
	links       map[string]ShortLinkData
	storageFile *os.File
}

type ShortLinkData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"orig_url"`
}

func NewFileShortLinkRepo(fileStoragePath string) (*FileShortLinkRepo, error) {
	repo := new(FileShortLinkRepo)
	file, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, constants.FileRWPerm)
	if err != nil {
		return nil, fmt.Errorf("failed open file storage: %w", err)
	}
	repo.storageFile = file
	links, err := repo.loadFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed load from file storage: %w", err)
	}

	repo.links = make(map[string]ShortLinkData, len(links))
	for _, link := range links {
		repo.links[link.ShortURL] = link
	}

	return repo, nil
}

func (repo *FileShortLinkRepo) CreateShortLink(shortID string, fullURL string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed create random: %w", err)
	}

	data := ShortLinkData{UUID: id.String(), ShortURL: shortID, OriginalURL: fullURL}
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

func (repo *FileShortLinkRepo) loadFromFile() ([]ShortLinkData, error) {
	scanner := bufio.NewScanner(repo.storageFile)
	var dataList []ShortLinkData

	for scanner.Scan() {
		dataBytes := scanner.Bytes()
		data := ShortLinkData{}
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

func (repo *FileShortLinkRepo) writeLink(link ShortLinkData) error {
	data, err := json.Marshal(link)
	if err != nil {
		return fmt.Errorf("failed serialize ShortLinkData: %w", err)
	}

	data = append(data, '\n')
	_, err = repo.storageFile.Write(data)
	return fmt.Errorf("failed write to file storage: %w", err)
}
