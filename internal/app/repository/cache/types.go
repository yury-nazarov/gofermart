package cache

type UserSessionInterface interface {
	// Add добавить токен в кеш. На пример LogIn
	Add(token string, userID int) error
	// GetUserIDByToken получить токен из кеша на пример для проверки сессии
	GetUserIDByToken(token string) (userID int, err error)
}
