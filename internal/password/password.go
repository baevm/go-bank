package password

import "golang.org/x/crypto/bcrypt"

func Hash(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), err
}

func Check(hash string, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
