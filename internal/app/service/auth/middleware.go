package auth

import (
	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
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
func HTTPTokenExist(auth UserInterface, logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из хедера и проверяем его в сессии
			token := r.Header.Get("Authorization")
			var user = models.UserDB{
				Token: token,
			}
			user, err := auth.IsSignIn(user)
			// Что то пошло не так с кешем
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			// Токен есть в кеше
			if user.ID != 0 {
				next.ServeHTTP(w, r)
			}

			// Остальные кейсы считаем пользователя не авторизованым
			w.WriteHeader(http.StatusUnauthorized)
		}
		return http.HandlerFunc(fn)
	}
}
