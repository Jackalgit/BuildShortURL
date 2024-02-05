package initialization

import (
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/database"
	dicturl "github.com/Jackalgit/BuildShortURL/internal/dictURL"
	"github.com/Jackalgit/BuildShortURL/internal/handlers"
	"github.com/Jackalgit/BuildShortURL/internal/userid"
)

func InitStorage() *handlers.ShortURL {

	if config.Config.DatabaseDSN != "" {
		return &handlers.ShortURL{
			Storage:         database.NewDataBase(),
			DictUserIDToken: userid.NewDictUserIDToken(),
		}
	}

	return &handlers.ShortURL{
		Storage:         dicturl.NewDictURL(),
		DictUserIDToken: userid.NewDictUserIDToken(),
	}

}
