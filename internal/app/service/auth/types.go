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
	// IsSignIn Проверяет авторизован ли пользователь
	IsSignIn(user models.UserDB) (models.UserDB, error)
}
