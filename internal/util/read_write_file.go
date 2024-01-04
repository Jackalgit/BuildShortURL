package util

import (
	"encoding/json"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"os"
)

func SaveURLToJsonFile(path string, originalURL string, shortURLKey string) error {

	if path == "" {
		return nil
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("[OpenFile] Не удалось открыть json file: %q\n", err)
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
		return fmt.Errorf("[Write to File] Не удалось записать LastURL json file: %q\n", err)
	}

	return nil
}

//func ReadURLFromFile(path string) error {
//
//	file, err := os.Open(path)
//	if err != nil {
//		return fmt.Errorf("[Open] Не удалось открыть json file: %q\n", err)
//	}
//	LastSku := PointSku{}
//	err = json.NewDecoder(file).Decode(&LastSku)
//	if err != nil {
//		return fmt.Errorf("[DecodeFile] Не удалось декодировать lastPointSeller.json: %q\n", err)
//	}
//
//	return nil
//}
