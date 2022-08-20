package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"github.com/yury-nazarov/gofermart/pkg/tools"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
)

type authLocalStruct struct {
	db           pg.DBInterface
	loginSession cache.UserSessionInterface
	logger       *log.Logger
}

func NewAuth(db pg.DBInterface, loginSession cache.UserSessionInterface, logger *log.Logger) authLocalStruct {
	return authLocalStruct{
		db:           db,
		loginSession: loginSession,
		logger:       logger,
	}
}

// SignUp регистрация пользователя
func (a authLocalStruct) SignUp(ctx context.Context, login string, password string) (token string, err400 error, err500 error) {
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
	hashPwd := hashPassword(password)

	// Записываем логин и хеш пароля в БД
	userID, err := a.db.NewUser(ctx, login, hashPwd)
	if err != nil {
		errString := fmt.Sprintf("can't create new user. err: %s", err)
		a.logger.Print(errString)
		return "", fmt.Errorf("%s", errString), nil
	}

	// Генерим Токен, добавляем токен и userID в сессию
	token = newToken()
	err = a.loginSession.Add(token, userID)
	if err != nil {
		errString := fmt.Sprintf("can't init user session. err: %s", err)
		a.logger.Print(errString)
		return "", nil, fmt.Errorf("%s", errString)
	}

	// Возвращаем токен для записи в заголовок
	return token, nil, nil
}

// SignIn вход пользователя
func (a authLocalStruct) SignIn(ctx context.Context, login string, password string) (token string, err error) {
	// Считаем хеш пароля
	hashPwd := hashPassword(password)

	// Проверяем в БД наличие пользователя
	userID, err := a.db.UserIsValid(ctx, login, hashPwd)
	if err != nil {
		errString := fmt.Sprintf("incorrect login or password: %s", err)
		a.logger.Print(errString)
		return "", tools.NewError401(errString)
	}

	// Генерим токен, добавляем в сессию
	token = newToken()
	err = a.loginSession.Add(token, userID)
	if err != nil {
		errString := fmt.Sprintf("error add token to session: %s", err)
		a.logger.Print(errString)
		return "", tools.NewError500(errString)
	}
	return token, nil
}

func (a authLocalStruct) IsUserSignIn(token string) (userID int, err500 error) {
	userID, err := a.loginSession.GetUserIDByToken(token)
	// Ошибка работы с кешем
	if err != nil {
		a.logger.Printf("session storage have error: %s", err500)
		return 0, err500
	}
	// Токен не найден
	if userID == 0 {
		a.logger.Printf("Token not exist")
		return 0, nil
	}
	// Токен найден
	return userID, nil
}

// newToken создает токен
func newToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// hashPassword считает хеш из пароля
func hashPassword(password string) string {
	hashPwd := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hashPwd)
}
