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

//[
//{
//"correlation_id": "<строковый идентификатор>",
//"original_url": "<URL для сокращения>"
//},
//...
//]
//[
//{
//"correlation_id": "<строковый идентификатор из объекта запроса>",
//"short_url": "<результирующий сокращённый URL>"
//},
//...
//]
