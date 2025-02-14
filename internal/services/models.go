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

func NewShortedLink(uuid string, corlId string, origURL string, url string, isDupl bool) *ShortedLink {
	return &ShortedLink{
		UUID:         uuid,
		CorelationID: corlId,
		OriginalURL:  origURL,
		URL:          url,
		IsDuplicated: isDupl,
	}
}
