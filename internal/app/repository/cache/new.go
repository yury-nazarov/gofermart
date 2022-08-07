package cache

type userSessionStruct struct {
	// map[Token]UserID
	data map[string]int
}


// NewLoginCache - Создает объект для хранения в RAM залогиненых пользователей и токенов для них
func NewLoginCache() *userSessionStruct {
	return &userSessionStruct{
		data: map[string]int{},
	}
}

// Add - 	добавить токен в кеш. На пример LogIn
func (l *userSessionStruct) Add(token string, userID int) error {
	l.data[token] = userID
	return nil
}

// GetUserIDByToken - 	получить id пользователя по токену
// 						error - пока заглушка, но если понадобится вынести во внешний кеш - то понадобится.
func (l *userSessionStruct) GetUserIDByToken(token string) (userID int, err error) {
	userID, ok := l.data[token]
	if ok {
		return userID, nil
	}
	return 0, nil
}