package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func CompareHashAndPassword(hashedPass string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), password)
	return err == nil
}
