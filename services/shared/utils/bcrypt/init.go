package bcrypt

import "golang.org/x/crypto/bcrypt"

type BCrypt struct {
	cost int
}

func Init(cost int) *BCrypt {
	return &BCrypt{cost: cost}
}

func (b *BCrypt) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BCrypt) Validate(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
