package authsdk

import (
	"github.com/neghi-go/auth-sdk/users"
)

type Auth struct {
	user *users.User
}

func New() (*Auth, error) {
	return &Auth{}, nil
}
