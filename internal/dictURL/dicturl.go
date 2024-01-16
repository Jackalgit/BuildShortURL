package dicturl

import "context"

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
