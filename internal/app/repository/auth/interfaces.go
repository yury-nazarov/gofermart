package auth

import "context"

/*
	Регистрации, аутентификации и проверки пользовательских сессий
*/

type UserInterface interface {
	// SignUp Регистрирует пользователя
	SignUp(ctx context.Context, login string, password string) (string, error, error)
	// SignIn Авторизует пользователя
	SignIn(ctx context.Context, login string, password string) (string, error)
	// IsUserSignIn Проверяет авторизован ли пользователь
	IsUserSignIn(token string) (bool, error)
}