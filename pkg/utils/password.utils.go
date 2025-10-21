package utils

import (
	"bme/pkg/errorext"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", errorext.New(err, errorext.ErrGeneralOccurrence)
	}

	return string(hash), nil
}

func VerifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return errorext.NewBadRequest(err, errorext.ErrInvalidPassword)
	}

	return nil
}
