package cache

type UserSessionInterface interface {
	// Add добавить токен в кеш
	Add(token string, userId int) error
	// GetUserIdByToken получить токен из кеша
	GetUserIdByToken(token string) (int, error)
}
