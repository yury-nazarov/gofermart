package auth

import (
	"log"
	"net/http"
)

/*
	Middleware для проверки наличия токена
	Если токен присланый в запросе есть в мапе
	map[string]string {
		Token: userID
	}
	то пользователь считается авторизованым
	id пользователя необходимо предать в запросе дальше, для аутентификации пользовательских операций

*/

// HTTPTokenExist Проверяет наличие токена
func HTTPTokenExist(user UserInterface, logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из хедера и проверяем его в сессии
			token := r.Header.Get("Authorization")
			userID, err500 := user.IsUserSignIn(token)
			// Что то пошло не так с кешем
			if err500 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Токен есть в кеше
			if userID != 0 {
				next.ServeHTTP(w, r)
			}

			// Остальные кейсы считаем пользователя не авторизованым
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		return http.HandlerFunc(fn)
	}
}
