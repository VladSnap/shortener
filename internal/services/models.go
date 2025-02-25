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
	IsDeleted    bool
}

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

type DeleteShortID struct {
	ShortURL string
	UserID   string
}

func NewDeleteShortID(shortURL string, userID string) DeleteShortID {
	return DeleteShortID{shortURL, userID}
}
