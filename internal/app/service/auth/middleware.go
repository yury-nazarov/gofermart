package auth

/*
	Middleware для проверки наличия токена
	Если токен присланый в запросе есть в мапе
	map[string]string {
		Token: username
	}
	то пользователь считается авторизованым

*/
