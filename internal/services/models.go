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
}
