package userID

import (
	"github.com/google/uuid"
)

type DictUserIdToken map[uuid.UUID]string

func NewDictUserIdToken() DictUserIdToken {
	return make(DictUserIdToken)
}

func (d DictUserIdToken) AddUserId(id uuid.UUID, tokenString string) {

	d[id] = tokenString

}
