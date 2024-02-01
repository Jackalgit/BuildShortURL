package dicturl

import (
	"context"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/Jackalgit/BuildShortURL/internal/util"
	"github.com/google/uuid"
)

type DictURL map[uuid.UUID]map[string][]byte

func NewDictURL() DictURL {
	return make(DictURL)
}

func (d DictURL) AddURL(ctx context.Context, userId uuid.UUID, shortURLKey string, originalURL []byte) error {

	userDictURL, foundDictUser := d[userId]
	if !foundDictUser {
		d[userId] = make(map[string][]byte)
		d[userId][shortURLKey] = originalURL
	} else {
		for key, value := range userDictURL {
			if string(value) == string(originalURL) {
				AddURLError := models.NewAddURLError(key)

				return AddURLError
			}
		}
		userDictURL[shortURLKey] = originalURL
	}

	util.SaveURLToJSONFile(config.Config.FileStoragePath, string(originalURL), shortURLKey)

	return nil

}

func (d DictURL) GetURL(ctx context.Context, userId uuid.UUID, shortURLKey string) ([]byte, bool) {

	userDictURL, foundDictUser := d[userId]
	if !foundDictUser {
		return nil, foundDictUser
	}

	originalURL, foundShortURLKey := userDictURL[shortURLKey]

	return originalURL, foundShortURLKey

}

func (d DictURL) AddBatchURL(ctx context.Context, userId uuid.UUID, batchList []models.BatchURL) error {

	userDictURL, foundDictUser := d[userId]
	if !foundDictUser {
		for _, v := range batchList {
			d[userId] = make(map[string][]byte)
			d[userId][v.ShortURL] = []byte(v.OriginalURL)
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

func (d DictURL) UserURLList(ctx context.Context, userId uuid.UUID) ([]models.ResponseUserURL, bool) {

	userDictURL, foundDictUser := d[userId]
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
