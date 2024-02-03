package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type RequestBatch struct {
	Correlation string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type ResponseBatch struct {
	Correlation string `json:"correlation_id"`
	ShortURL    string `json:"short_url"`
}

type ResponseUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type DeleteShortURL struct {
	ShortURL string `json:"short_url"`
}
