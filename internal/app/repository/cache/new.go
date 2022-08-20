package cache

import (
	"log"
	"sync"
)

type userSessionStruct struct {
	// map[Token]UserID
	data map[string]int
	mu   sync.RWMutex
}

// NewLoginCache - Создает объект для хранения в RAM залогиненых пользователей и токенов для них
func NewLoginCache() *userSessionStruct {
	return &userSessionStruct{
		data: map[string]int{},
	}
}

// Add - добавить токен в кеш. На пример LogIn
func (l *userSessionStruct) Add(token string, userID int) error {
	// Берем мутекс на момент записи
	l.mu.Lock()
	defer l.mu.Unlock()

	l.data[token] = userID

	return nil
}

// GetUserIDByToken - 	получить id пользователя по токену
// 						error - пока заглушка, но если понадобится вынести во внешний кеш - то понадобится.
func (l *userSessionStruct) GetUserIDByToken(token string) (userID int, err error) {
	log.Println(l.data)
	userID, ok := l.data[token]
	if ok {
		return userID, nil
	}
	return 0, nil
}
