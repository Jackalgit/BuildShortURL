package initialization

import (
	"context"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/database"
	dicturl "github.com/Jackalgit/BuildShortURL/internal/dictURL"
	"github.com/Jackalgit/BuildShortURL/internal/filestorage"
	"github.com/Jackalgit/BuildShortURL/internal/handlers"
)

func InitStorage(ctx context.Context) *handlers.ShortURL {

	if config.Config.DatabaseDSN != "" {
		return &handlers.ShortURL{
			Ctx:     ctx,
			Storage: database.NewDataBase(),
		}
	}

	if config.Config.FileStoragePath != "" {
		return &handlers.ShortURL{
			Ctx:     ctx,
			Storage: filestorage.NewFileStorage(config.Config.FileStoragePath),
		}
	}

	return &handlers.ShortURL{
		Ctx:     ctx,
		Storage: dicturl.NewDictURL(),
	}

}