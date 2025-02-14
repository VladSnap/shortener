package data

type ShortLinkData struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"orig_url" db:"orig_url"`
}

func NewShortLinkData(id string, shortURL string, origURL string) *ShortLinkData {
	return &ShortLinkData{
		UUID:        id,
		ShortURL:    shortURL,
		OriginalURL: origURL,
	}
}
