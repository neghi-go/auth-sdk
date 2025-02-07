package users

import "github.com/neghi-go/auth-sdk/utils"

type UserModelCreateOptions func(*UserModel) error

func SetPassword(password string) UserModelCreateOptions {
	hash := utils.NewHasher()
	return func(um *UserModel) error {
		salt, err := utils.GenerateSalt(64)
		if err != nil {
			return err
		}
		hashedPassword := hash.Hash(password, salt)
		um.EncryptedPasswordSalt = salt
		um.EncryptedPassword = hashedPassword
		return nil
	}
}

func SetAudience(aud string) UserModelCreateOptions {
	return func(um *UserModel) error {
		um.Aud = aud
		return nil
	}
}

type UserModelUpdateOptions func(*UserModel) error
