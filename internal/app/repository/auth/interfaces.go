package auth


/*
	Регистрации, аутентификации и проверки пользовательских сессий
*/

type UserInterface interface {
	// SignUp Регистрирует пользователя
	SignUp(login string, password string) error
	// SignIn Авторизует пользователя
	SignIn(login string, password string) error
	// IsUserSignIn Проверяет авторизован ли пользователь
	IsUserSignIn(token string) (bool, error)
}