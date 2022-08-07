package cache

import "fmt"

type userSessionStruct struct {
	// map[Token]UserID
	data map[string]int
}


// NewLoginCache Создает объект для хранения в RAM залогиненых пользователей и токенов для них
func NewLoginCache() *userSessionStruct {
	return &userSessionStruct{
		data: map[string]int{},
	}
}


func (l *userSessionStruct) Add(token string, userID int) error {
	l.data[token] = userID
	return nil
}


func (l *userSessionStruct) GetUserIDByToken(token string) (int, error) {
	userID, ok := l.data[token]
	if ok {
		return userID, nil
	}
	return 0, fmt.Errorf("token not exist in session")
}