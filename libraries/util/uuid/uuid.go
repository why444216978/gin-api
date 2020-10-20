package go_util

import (
	"github.com/google/uuid"
)

func GetUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
