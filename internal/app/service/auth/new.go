package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"github.com/yury-nazarov/gofermart/pkg/tools"
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
func (a authLocalStruct) SignUp(ctx context.Context, user models.UserDB) (models.UserDB, error) {
	ok, err := a.db.UserExist(ctx, user.Login)
	// Error500
	if err != nil {
		errString := fmt.Sprintf("UserExist sql querry error: %s", err)
		a.logger.Print(errString)
		return user, tools.NewError500(errString)
	}
	// Error400
	if ok {
		errString := fmt.Sprintf("user exist: %s", err)
		a.logger.Print(errString)
		return user, tools.NewError400(errString)
	}

	// Считаем хеш пароля и перезаписываем переменную
	user.Password = hashPassword(user.Password)

	// Записываем логин и хеш пароля в БД
	//userID, err := a.db.NewUser(ctx, user)
	user, err = a.db.NewUser(ctx, user)
	if err != nil {
		errString := fmt.Sprintf("can't create new user. err: %s", err)
		a.logger.Print(errString)
		return user, tools.NewError400(errString)
	}

	// Генерим Токен, добавляем токен и userID в сессию
	user.Token = newToken()
	err = a.loginSession.Add(user)
	if err != nil {
		errString := fmt.Sprintf("can't init user session. err: %s", err)
		a.logger.Print(errString)
		return user, tools.NewError500(errString)
	}

	// Возвращаем токен для записи в заголовок
	return user, nil
}

// SignIn вход пользователя
func (a authLocalStruct) SignIn(ctx context.Context, user models.UserDB) (models.UserDB,error) {
	// Считаем хеш пароля
	user.Password = hashPassword(user.Password)

	// Проверяем в БД наличие пользователя
	user, err := a.db.UserIsValid(ctx, user)
	if err != nil {
		errString := fmt.Sprintf("incorrect login or password: %s", err)
		a.logger.Print(errString)
		return user, tools.NewError401(errString)
	}

	// Генерим токен, добавляем в сессию
	user.Token = newToken()
	err = a.loginSession.Add(user)
	if err != nil {
		errString := fmt.Sprintf("error add token to session: %s", err)
		a.logger.Print(errString)
		return user, tools.NewError500(errString)
	}
	return user, nil
}

//func (a authLocalStruct) IsSignIn(token string) (userID int, err error) {
func (a authLocalStruct) IsSignIn(user models.UserDB) (models.UserDB, error) {
	userID, err := a.loginSession.GetUserIDByToken(user.Token)
	// Ошибка работы с кешем
	if err != nil {
		a.logger.Printf("session storage have error: %s", err)
		return user, err
	}
	// Токен не найден
	if userID == 0 {
		a.logger.Printf("Token not exist")
		return user, nil
	}
	// Токен найден
	user.ID = userID
	return user, nil
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
