package models

type BatchList struct {
	List []BatchURL
}

type BatchURL struct {
	Correlation string `json:"correlation_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
