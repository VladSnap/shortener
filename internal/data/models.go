package data

type ShortLinkData struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"orig_url" db:"orig_url"`
	UserID      string `json:"user_id" db:"user_id"`
	IsDeleted   bool   `json:"is_deleted" db:"is_deleted"`
}

func NewShortLinkData(id string, shortURL string, origURL string, userID string) *ShortLinkData {
	return &ShortLinkData{
		UUID:        id,
		ShortURL:    shortURL,
		OriginalURL: origURL,
		UserID:      userID,
	}
}

type DeleteShortData struct {
	ShortURL string
	UserID   string
}

func NewDeleteShortData(shortURL string, userID string) DeleteShortData {
	return DeleteShortData{shortURL, userID}
}
