package utils

import (
	"github.com/alexedwards/argon2id"
)

func GetHashedPassword(rawPassword string, params *argon2id.Params) (string, error) {
	hashedPassword, err := argon2id.CreateHash(rawPassword, params)
	return hashedPassword, err
}

func CompareWithHashedPassword(rawPassword, hashedPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(rawPassword, hashedPassword)
	return match, err
}
