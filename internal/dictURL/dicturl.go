package dicturl

type DictURL map[string][]byte

func NewDictURL() DictURL {
	return make(DictURL)
}

func (d DictURL) AddURL(key string, originalURL []byte) {
	d[key] = originalURL

}

func (d DictURL) GetURL(key string) ([]byte, bool) {
	originalURL, found := d[key]

	return originalURL, found

}
