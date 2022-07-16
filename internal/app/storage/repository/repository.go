package repository

import "context"

// Пакет содержит общие интерфейсы и структуры для работы с БД

type Repository interface {
	Ping() bool
	GetToken(ctx context.Context, token string) (bool, error)
}