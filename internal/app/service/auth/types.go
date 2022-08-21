package auth

import (
	"context"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
)

type UserInterface interface {
	// SignUp Регистрирует пользователя
	SignUp(ctx context.Context, user models.UserDB) (token string, err error)
	// SignIn Авторизует пользователя
	SignIn(ctx context.Context, user models.UserDB) (token string, err error)
	// IsUserSignIn Проверяет авторизован ли пользователь
	IsUserSignIn(token string) (userID int, err500 error)
}
