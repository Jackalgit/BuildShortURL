package filestorage

import (
	"encoding/json"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"os"
)

func SaveURLToJSONFile(path string, originalURL string, shortURLKey string) error {

	if path == "" {
		return nil
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("[OpenFile] Не удалось открыть json file: %q", err)

	}
	defer file.Close()

	LastURL := models.FileStorageDictURL{
		OriginalURL: originalURL,
		ShortURL:    shortURLKey,
	}

	data, _ := json.MarshalIndent(&LastURL, "", " ")
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("[Write to File] Не удалось записать LastURL json file: %q", err)

	}

	return nil
}

func SaveListURLToJSONFile(path string, batchList []models.BatchURL) error {

	if path == "" {
		return nil
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("[OpenFile] Не удалось открыть json file: %q", err)

	}
	defer file.Close()

	data, _ := json.MarshalIndent(batchList, "", " ")
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("[Write to File] Не удалось записать batchList json file: %q", err)

	}

	return nil

}
