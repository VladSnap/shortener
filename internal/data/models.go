package data

type ShortLinkData struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"orig_url" db:"orig_url"`
}
