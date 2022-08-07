package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository"
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
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
		errString := fmt.Sprintf("NewUser sql querry error: %s", err)
		a.logger.Print(errString)
		return "", fmt.Errorf("%s", errString), nil
	}

	// Генерим Токен, добавляем токен и userID в сессию
	token = newToken()
	err = a.loginSession.Add(token, userID)
	if err != nil {
		errString := fmt.Sprintf("add token to session: %s", err)
		a.logger.Print(errString)
		return "", nil, fmt.Errorf("%s", errString)
	}

	// Возвращаем токен для записи в заголовок
	a.logger.Printf("token: %s success", token)
	return token, nil, nil
}

// SignIn вход пользователя
func (a authLocalStruct) SignIn(ctx context.Context, login string, password string) (token string, err401 error, err500 error) {
	// Считаем хеш пароля
	hashPwd := hashPassword(password)
	// TODO: DEBUG
	a.logger.Printf("DEBUG: Calculate hash. User %s, has: %s", login, hashPwd)

	// Проверяем в БД наличие пользователя
	userID, err401 := a.db.UserIsValid(ctx, login, hashPwd)
	if err401 != nil {
		errString := fmt.Sprintf("incorrect login or password: %s", err401)
		a.logger.Print(errString)
		return "", fmt.Errorf("%s", errString), nil
	}
	// TODO: DEBUG
	a.logger.Printf("DEBUG: User %s exist with hash pwd: %s", login, hashPwd)

	// Генерим токен, добавляем в сессию
	token = newToken()
	err := a.loginSession.Add(token, userID)
	if err != nil {
		errString := fmt.Sprintf("error add token to session: %s", err)
		a.logger.Print(errString)
		return "", nil, fmt.Errorf("%s", errString)
	}
	// TODO: DEBUG
	a.logger.Printf("DEBUG: User %s: Add  user id: %d and token: %s in to session", login, userID, token)

	return token, nil, nil

}

func (a authLocalStruct) IsUserSignIn(token string) (userID int, err500 error) {
	userID, err := a.loginSession.GetUserIDByToken(token)
	// Ошибка работы с кешем
	if err != nil {
		a.logger.Printf("session storage have error: %s", err500)
		return 0,  err500
	}
	// Токен не найден
	if userID == 0 {
		a.logger.Printf("Token not exist")
		return 0, nil
	}
	// Токен найден
	return userID,  nil
}

// newToken создает токен
func newToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//// hashPassword считает хеш из пароля - при этом каждый раз разный. Больше подходит для токена))
//func hashPassword(password string) (string, error) {
//	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
//	return string(bytes), err
//}

// hashPassword считает хеш из пароля
func hashPassword(password string) string {
	hashPwd := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hashPwd)
}