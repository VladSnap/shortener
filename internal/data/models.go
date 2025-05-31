// Package data хранит модели данных для БД.
package data

// ShortLinkData - Структура таблицы БД сокращенной ссылки.
type ShortLinkData struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"orig_url" db:"orig_url"`
	UserID      string `json:"user_id" db:"user_id"`
	IsDeleted   bool   `json:"is_deleted" db:"is_deleted"`
}

// NewShortLinkData - Создает новую структуру ShortLinkData с указателем.
func NewShortLinkData(id string, shortURL string, origURL string, userID string) *ShortLinkData {
	return &ShortLinkData{
		UUID:        id,
		ShortURL:    shortURL,
		OriginalURL: origURL,
		UserID:      userID,
	}
}

// DeleteShortData - Структура запроса для удаления сокращенной ссылки.
type DeleteShortData struct {
	ShortURL string
	UserID   string
}

// NewDeleteShortData - Создает новую структуру DeleteShortData с указателем.
func NewDeleteShortData(shortURL string, userID string) DeleteShortData {
	return DeleteShortData{shortURL, userID}
}

// StatsData - Статистика по пользователям и всем сокращенным ссылкам.
type StatsData struct {
	Urls  int
	Users int
}

// NewStatsData - Создает новую структуру StatsData с указателем.
func NewStatsData(urls, users int) *StatsData {
	return &StatsData{Urls: urls, Users: users}
}
