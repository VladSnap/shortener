// Package services этот пакет хранит модели и сервисы суровня бизнес логики.
package services

// OriginalLink - Структура и доменный объект оригинальной ссылки.
type OriginalLink struct {
	CorelationID string
	URL          string
}

// ShortedLink - Структура и доменный объект сокращенной ссылки.
type ShortedLink struct {
	UUID         string
	CorelationID string
	OriginalURL  string
	URL          string
	IsDuplicated bool
	IsDeleted    bool
}

// NewShortedLink - Создает новую структуру ShortedLink с указателем.
func NewShortedLink(uuid string, corlID string, origURL string, url string, isDupl bool, isDel bool) *ShortedLink {
	return &ShortedLink{
		UUID:         uuid,
		CorelationID: corlID,
		OriginalURL:  origURL,
		URL:          url,
		IsDuplicated: isDupl,
		IsDeleted:    isDel,
	}
}

// DeleteShortID - Структура запроса удаления сокращенной ссылки.
type DeleteShortID struct {
	ShortURL string
	UserID   string
}

// NewDeleteShortID - Создает новую структуру DeleteShortID с указателем.
func NewDeleteShortID(shortURL string, userID string) DeleteShortID {
	return DeleteShortID{shortURL, userID}
}
