package uuid

import (
	"github.com/google/uuid"
)

func NewUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
