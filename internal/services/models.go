package services

type OriginalLink struct {
	CorelationID string
	URL          string
}

type ShortedLink struct {
	UUID         string
	CorelationID string
	OriginalURL  string
	URL          string
	IsDuplicated bool
}

func NewShortedLink(uuid string, corlID string, origURL string, url string, isDupl bool) *ShortedLink {
	return &ShortedLink{
		UUID:         uuid,
		CorelationID: corlID,
		OriginalURL:  origURL,
		URL:          url,
		IsDuplicated: isDupl,
	}
}
