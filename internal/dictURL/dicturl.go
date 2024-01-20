package dicturl

import (
	"context"
	"github.com/Jackalgit/BuildShortURL/internal/models"
)

type DictURL map[string][]byte

func NewDictURL() DictURL {
	return make(DictURL)
}

func (d DictURL) AddURL(ctx context.Context, shortURLKey string, originalURL []byte) {
	d[shortURLKey] = originalURL

}

func (d DictURL) GetURL(ctx context.Context, shortURLKey string) ([]byte, bool) {
	originalURL, found := d[shortURLKey]

	return originalURL, found

}

func (d DictURL) AddBatchURL(ctx context.Context, batchList *models.BatchList) {
	for _, v := range batchList.List {
		d[v.ShortURL] = []byte(v.OriginalURL)
	}

}
