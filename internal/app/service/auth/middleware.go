package auth

import (
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
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
func HTTPTokenExist(session cache.UserSessionInterface, logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из хедера
			token := r.Header.Get("Authorization")

			// Если токен есть в кеше - считаем пользователя авторизованым
			userID, err := session.GetUserIDByToken(token)
			if err != nil {
				logger.Printf("session storage have error: %s", err)
			}
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


//func HTTPTokenExist(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token := r.Header.Get("Authorization")
//		// Если токен не установлен - считаем пользователя не авторизованым
//		if len(token) == 0 {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		tokenExist, err :=
//		if {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//		next.ServeHTTP(w, r)
//
//	})
//}
