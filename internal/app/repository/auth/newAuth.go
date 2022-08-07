package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository"
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"

	"golang.org/x/crypto/bcrypt"
)

type authLocalStruct struct {
	db repository.DBInterface
	loginSession cache.UserSessionInterface
	logger *log.Logger
}

func NewAuth(db repository.DBInterface, loginSession cache.UserSessionInterface, logger *log.Logger) UserInterface{
	return authLocalStruct{
		db: db,
		loginSession: loginSession,
		logger: logger,
	}
}

// SignUp регистрация пользователя
func (a authLocalStruct) SignUp(ctx context.Context, login string, password string) (string, error, error) {
	ok, err := a.db.UserExist(ctx, login)
	// Error500
	if err != nil {
		errString := fmt.Sprintf("UserExist sql querry error: %s", err)
		a.logger.Print(errString)
		return "", nil, fmt.Errorf("%s", errString)
	}
	// Error409
	if ok {
		errString := fmt.Sprintf("user exist: %s", err)
		a.logger.Print(errString)
		return "", fmt.Errorf("%s", errString), nil
	}

	// Считаем хеш пароля
	hashPwd, err := hashPassword(password)
	if err != nil {
		errString := fmt.Sprintf("hashPassword error: %s", err)
		a.logger.Print(errString)
		return "", nil, fmt.Errorf("%s", errString)
	}
	// Записываем логин и хеш пароля в БД
	userId, err := a.db.NewUser(ctx, login, hashPwd)
	if err != nil {
		errString := fmt.Sprintf("NewUser sql querry error: %s", err)
		a.logger.Print(errString)
		return "", fmt.Errorf("%s", errString), nil
	}

	// Генерим Токен, добавляем токен и userId в сессию
	token := newToken()
	err = a.loginSession.Add(token, userId)
	if err != nil {
		errString := fmt.Sprintf("add token to session: %s", err)
		a.logger.Print(errString)
		return "", nil, fmt.Errorf("%s", errString)
	}

	// Возвращаем токен для записи в заголовок
	a.logger.Printf("token: %s success", token)
	return token, nil, nil
}


func (a authLocalStruct) SignIn(ctx context.Context, login string, password string) (string, error) {
	return "", nil
}

func (a authLocalStruct) IsUserSignIn(token string) (bool, error) {
	return true, nil
}

// newToken создает токен
func newToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// hashPassword считает хеш из пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}