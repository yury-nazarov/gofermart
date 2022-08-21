package auth

import (
	"context"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
)

type UserInterface interface {
	// SignUp Регистрирует пользователя
	SignUp(ctx context.Context, user models.UserDB) (models.UserDB, error)
	// SignIn Авторизует пользователя
	SignIn(ctx context.Context, user models.UserDB) (models.UserDB, error)
	// IsUserSignIn Проверяет авторизован ли пользователь
	IsUserSignIn(token string) (userID int, err500 error)
}
