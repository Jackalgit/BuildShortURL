package userid

import (
	"github.com/google/uuid"
)

type DictUserIDToken map[uuid.UUID]string

func NewDictUserIDToken() DictUserIDToken {
	return make(DictUserIDToken)
}

func (d DictUserIDToken) AddUserID(id uuid.UUID, tokenString string) {

	d[id] = tokenString

}
