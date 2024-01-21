package dicturl

import (
	"context"
	"github.com/Jackalgit/BuildShortURL/internal/models"
)

type DictURL map[string][]byte

func NewDictURL() DictURL {
	return make(DictURL)
}

func (d DictURL) AddURL(ctx context.Context, shortURLKey string, originalURL []byte) error {
	for key, value := range d {
		if string(value) == string(originalURL) {
			AddURLError := models.NewAddURLError(key)

			return AddURLError
		}
	}

	d[shortURLKey] = originalURL
	return nil

}

func (d DictURL) GetURL(ctx context.Context, shortURLKey string) ([]byte, bool) {
	originalURL, found := d[shortURLKey]

	return originalURL, found

}

func (d DictURL) AddAPIShortURL(ctx context.Context, shortURLKey string, originalURL []byte) {

	d[shortURLKey] = originalURL

	return

}

func (d DictURL) AddBatchURL(ctx context.Context, batchList []models.BatchURL) error {
	for _, v := range batchList {
		for key, value := range d {
			if string(value) == v.OriginalURL {
				AddURLError := models.NewAddURLError(key)

				return AddURLError
			}
		}

		d[v.ShortURL] = []byte(v.OriginalURL)
	}
	return nil

}
