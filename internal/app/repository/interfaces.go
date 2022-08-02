package repository

/*
	Описывает интерфейсы работы со слоем
	Работа с источниками данных
*/

import (
	"context"
)

// DBInterface методы работы с релиационными БД
type DBInterface interface {
	Ping() bool
	GetToken(ctx context.Context, token string) (bool, error)
	// GetUser NewUser: Операции с пользователями
	GetUser(ctx context.Context, login string) (bool,error)
	NewUser(ctx context.Context, login string, pwd string) (string,error)
}

// AccrualInterface методы работы с системой рассчета баллов
type AccrualInterface interface {

}
