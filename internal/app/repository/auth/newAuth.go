package auth

import (
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository"
)

type authLocalStruct struct {
	db repository.DBInterface
	logger *log.Logger
}

func NewAuth(db repository.DBInterface, logger *log.Logger) UserInterface{
	return authLocalStruct{
		db: db,
		logger: logger,
	}
}


func (a authLocalStruct) SignUp(login string, password string) error {
	return nil
}

func (a authLocalStruct) SignIn(login string, password string) error {
	return nil
}

func (a authLocalStruct) IsUserSignIn(token string) (bool, error) {
	return true, nil
}