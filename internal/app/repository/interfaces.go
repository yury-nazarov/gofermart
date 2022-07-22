package repository

import "context"

// DBInterface методы работы с релиационными БД
type DBInterface interface {
	Ping() bool
	GetToken(ctx context.Context, token string) (bool, error)
}

// AccrualInterface методы работы с системой рассчета баллов
type AccrualInterface interface {

}
