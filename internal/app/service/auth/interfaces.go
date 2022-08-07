package auth

import "context"

/*
	Регистрации, аутентификации и проверки пользовательских сессий
*/

type UserInterface interface {
	// SignUp Регистрирует пользователя
	SignUp(ctx context.Context, login string, password string) (token string, err400 error, err500 error)
	// SignIn Авторизует пользователя
	SignIn(ctx context.Context, login string, password string) (token string, err401 error, err500 error)
	// IsUserSignIn Проверяет авторизован ли пользователь
	IsUserSignIn(token string) (ok bool, err error)
}