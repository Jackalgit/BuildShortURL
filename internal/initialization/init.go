package initialization

import (
	"context"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/database"
	dicturl "github.com/Jackalgit/BuildShortURL/internal/dictURL"
	"github.com/Jackalgit/BuildShortURL/internal/handlers"
	"github.com/Jackalgit/BuildShortURL/internal/userID"
)

func InitStorage(ctx context.Context) *handlers.ShortURL {

	if config.Config.DatabaseDSN != "" {
		return &handlers.ShortURL{
			Ctx:             ctx,
			Storage:         database.NewDataBase(ctx),
			DictUserIdToken: userID.NewDictUserIdToken(),
		}
	}

	return &handlers.ShortURL{
		Ctx:             ctx,
		Storage:         dicturl.NewDictURL(),
		DictUserIdToken: userID.NewDictUserIdToken(),
	}

}
