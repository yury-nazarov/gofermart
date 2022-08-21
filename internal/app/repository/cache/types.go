package cache

import "github.com/yury-nazarov/gofermart/internal/app/repository/models"

type UserSessionInterface interface {
	// Add добавить токен в кеш. На пример LogIn
	Add(user models.UserDB) error
	//Add(token string, userID int) error
	// GetUserIDByToken получить токен из кеша на пример для проверки сессии
	GetUserIDByToken(token string) (userID int, err error)
}
