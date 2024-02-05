package dicturl

import (
	"context"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/filestorage"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/google/uuid"
)

type DictURL map[uuid.UUID]map[string][]byte

func NewDictURL() DictURL {
	return make(DictURL)
}

func (d DictURL) AddURL(ctx context.Context, userID uuid.UUID, shortURLKey string, originalURL []byte) error {

	userDictURL, foundDictUser := d[userID]
	if !foundDictUser {
		d[userID] = make(map[string][]byte)
		d[userID][shortURLKey] = originalURL
	} else {
		for key, value := range userDictURL {
			if string(value) == string(originalURL) {
				AddURLError := models.NewAddURLError(key)

				return AddURLError
			}
		}
		userDictURL[shortURLKey] = originalURL
	}

	filestorage.SaveURLToJSONFile(config.Config.FileStoragePath, string(originalURL), shortURLKey)

	return nil

}

func (d DictURL) GetURL(ctx context.Context, userID uuid.UUID, shortURLKey string) ([]byte, models.StatusURL) {

	userDictURL, foundDictUser := d[userID]
	if !foundDictUser {
		for _, dictUser := range d {
			for short, origin := range dictUser {
				if short == shortURLKey {
					return origin, models.NewStatusURL(true, false)
				}
			}
		}

		return nil, models.NewStatusURL(false, false)
	}

	origin, foundShortURLKey := userDictURL[shortURLKey]

	return origin, models.NewStatusURL(foundShortURLKey, false)

}

func (d DictURL) AddBatchURL(ctx context.Context, userID uuid.UUID, batchList []models.BatchURL) error {

	userDictURL, foundDictUser := d[userID]
	if !foundDictUser {
		for _, v := range batchList {
			d[userID] = make(map[string][]byte)
			d[userID][v.ShortURL] = []byte(v.OriginalURL)
		}
	} else {
		for _, v := range batchList {
			for key, value := range userDictURL {
				if string(value) == v.OriginalURL {
					AddURLError := models.NewAddURLError(key)

					return AddURLError
				}
			}
			userDictURL[v.ShortURL] = []byte(v.OriginalURL)
		}
	}

	filestorage.SaveListURLToJSONFile(config.Config.FileStoragePath, batchList)

	return nil

}

func (d DictURL) UserURLList(ctx context.Context, userID uuid.UUID) ([]models.ResponseUserURL, bool, error) {

	userDictURL, foundDictUser := d[userID]
	if !foundDictUser {
		return nil, foundDictUser, nil
	}

	var responseUserURLList []models.ResponseUserURL

	for k, v := range userDictURL {

		responseUserURL := models.ResponseUserURL{
			ShortURL:    fmt.Sprint(config.Config.BaseAddress, "/", k),
			OriginalURL: string(v),
		}

		responseUserURLList = append(responseUserURLList, responseUserURL)

	}

	return responseUserURLList, true, nil

}
