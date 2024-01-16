package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"os"
)

type FileStorage struct {
	Path string
}

func NewFileStorage(path string) FileStorage {
	return FileStorage{
		Path: path,
	}
}
func (f FileStorage) AddURL(ctx context.Context, shortURLKey string, originalURL []byte) {
	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Errorf("[OpenFile] Не удалось открыть json file: %q", err)

	}
	defer file.Close()

	LastURL := models.FileStorageDictURL{
		OriginalURL: string(originalURL),
		ShortURL:    shortURLKey,
	}

	data, _ := json.MarshalIndent(&LastURL, "", " ")
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		fmt.Errorf("[Write to File] Не удалось записать LastURL json file: %q", err)

	}

}

func (f FileStorage) GetURL(ctx context.Context, shortURLKey string) ([]byte, bool) {
	return nil, false
}
