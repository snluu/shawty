package data

import (
	"errors"
	"shawty/utils"
)

// ShortID constructs a short ID
func ShortID(id uint64, r string) string {
	return r + utils.ToSafeBase(id)
}

// FullID deconstructs a full ID
func FullID(str string) (uint64, string, error) {
	if len(str) < 2 {
		return 0, "", errors.New("Cannot deconstruct " + str)
	}
	return utils.ToDec(str[1:]), str[:1], nil
}
