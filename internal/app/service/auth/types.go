package auth

import "context"

type UserInterface interface {
	// SignUp Регистрирует пользователя
	SignUp(ctx context.Context, login string, password string) (token string, err error)
	// SignIn Авторизует пользователя
	SignIn(ctx context.Context, login string, password string) (token string, err error)
	// IsUserSignIn Проверяет авторизован ли пользователь
	IsUserSignIn(token string) (userID int, err500 error)
}

// User структура для JSON Unmarshal из HTTP Request
type User struct {
	Login    string `json:"login,omitempty" binding:"required"`
	Password string `json:"password" binding:"required"`
}
