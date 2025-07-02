package auth_service

import "github.com/3Eeeecho/go-gin-example/models"

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (bool, error) {
	return models.CheckUser(a.Username, a.Password)
}
