package data

import (
	"encoding/json"
	"os"

	"bufio"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/google/uuid"
)

type FileShortLinkRepo struct {
	links       map[string]ShortLinkData
	storageFile *os.File
}

type ShortLinkData struct {
	UUID        string
	ShortURL    string
	OriginalURL string
}

func NewFileShortLinkRepo(fileStoragePath string) (*FileShortLinkRepo, error) {
	repo := new(FileShortLinkRepo)
	file, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	repo.storageFile = file
	links, err := repo.loadFromFile()
	if err != nil {
		return nil, err
	}

	repo.links = make(map[string]ShortLinkData, len(links))
	for _, link := range links {
		repo.links[link.ShortURL] = link
	}

	return repo, nil
}

func (repo *FileShortLinkRepo) CreateShortLink(shortID string, fullURL string) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	data := ShortLinkData{UUID: uuid.String(), ShortURL: shortID, OriginalURL: fullURL}
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
		log.Zap.Error("File storage close error", err.Error())
	}
	log.Zap.Info("File storage closed")

	return err
}

func (repo *FileShortLinkRepo) loadFromFile() ([]ShortLinkData, error) {
	scanner := bufio.NewScanner(repo.storageFile)
	var dataList []ShortLinkData

	for scanner.Scan() {
		dataBytes := scanner.Bytes()
		data := ShortLinkData{}
		err := json.Unmarshal(dataBytes, &data)
		if err != nil {
			return nil, err
		}

		dataList = append(dataList, data)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dataList, nil
}

func (repo *FileShortLinkRepo) writeLink(link ShortLinkData) error {
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = repo.storageFile.Write(data)
	return err
}
