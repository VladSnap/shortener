package repos

import (
	"context"

	"github.com/VladSnap/shortener/internal/data"
)

// InMemoryShortLinkRepo - Репозиторий для доступа к хранилищу в оперативной памяти сокращателя ссылок.
type InMemoryShortLinkRepo struct {
	links map[string]*data.ShortLinkData
}

// NewShortLinkRepo - Создает новую структуру InMemoryShortLinkRepo с указателем.
func NewShortLinkRepo() *InMemoryShortLinkRepo {
	repo := new(InMemoryShortLinkRepo)
	repo.links = make(map[string]*data.ShortLinkData)
	return repo
}

// Add - Сохраняет структуру сокращенной ссылки в памяти.
func (repo *InMemoryShortLinkRepo) Add(ctx context.Context, link *data.ShortLinkData) (*data.ShortLinkData, error) {
	repo.links[link.ShortURL] = link
	return link, nil
}

// AddBatch - Сохраняет пачку структур сокращенных ссылок в памяти.
func (repo *InMemoryShortLinkRepo) AddBatch(ctx context.Context, links []*data.ShortLinkData) (
	[]*data.ShortLinkData, error) {
	for _, link := range links {
		repo.links[link.ShortURL] = link
	}
	return links, nil
}

// Get - Читает полную ссылку по сокращенной ссылке.
func (repo *InMemoryShortLinkRepo) Get(ctx context.Context, shortID string) (*data.ShortLinkData, error) {
	return repo.links[shortID], nil
}

// GetAllByUserID - Получить все сокращенные ссылки указанного пользователя.
func (repo *InMemoryShortLinkRepo) GetAllByUserID(ctx context.Context, userID string) (
	[]*data.ShortLinkData, error) {
	return make([]*data.ShortLinkData, 0), nil
}

// DeleteBatch - Удаляет пачку структур сокращенных ссылок из файла.
func (repo *InMemoryShortLinkRepo) DeleteBatch(ctx context.Context, shortIDs []data.DeleteShortData) error {
	for _, sid := range shortIDs {
		link := repo.links[sid.ShortURL]
		if link.UserID == sid.UserID {
			link.IsDeleted = true
		}
	}
	return nil
}

// GetStats - Получает статистику о пользователях и всех ссылках.
func (repo *InMemoryShortLinkRepo) GetStats(ctx context.Context) (*data.StatsData, error) {
	data := data.NewStatsData(len(repo.links), repo.calcAllUsers())
	return data, nil
}

func (repo *InMemoryShortLinkRepo) calcAllUsers() int {
	users := make(map[string]bool)

	for _, l := range repo.links {
		users[l.UserID] = true
	}

	return len(users)
}
