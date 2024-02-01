package dicturl

import (
	"context"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/Jackalgit/BuildShortURL/internal/util"
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
	fmt.Print(d)

	util.SaveURLToJSONFile(config.Config.FileStoragePath, string(originalURL), shortURLKey)

	return nil

}

func (d DictURL) GetURL(ctx context.Context, userID uuid.UUID, shortURLKey string) ([]byte, bool) {

	userDictURL, foundDictUser := d[userID]
	if !foundDictUser {
		for _, dictUser := range d {
			for short, origin := range dictUser {
				if short == shortURLKey {
					return origin, true
				}
			}
		}

		return nil, foundDictUser
	}

	origin, foundShortURLKey := userDictURL[shortURLKey]

	return origin, foundShortURLKey

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

	util.SaveListURLToJSONFile(config.Config.FileStoragePath, batchList)

	return nil

}

func (d DictURL) UserURLList(ctx context.Context, userID uuid.UUID) ([]models.ResponseUserURL, bool) {

	userDictURL, foundDictUser := d[userID]
	if !foundDictUser {
		return nil, foundDictUser
	}

	var responseUserURLList []models.ResponseUserURL

	for k, v := range userDictURL {

		responseUserURL := models.ResponseUserURL{
			ShortURL:    k,
			OriginalURL: string(v),
		}

		responseUserURLList = append(responseUserURLList, responseUserURL)

	}

	return responseUserURLList, true

}
