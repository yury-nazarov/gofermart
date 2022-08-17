package withdraw

import "context"

/*
	Типы для работы со слоем бизнес логики: вывод средств и запрос текущего баланса
*/

type BalanceInterface interface {
	CurrentBalance(ctx context.Context, userID int) (Balance, error)
}

// Balance Для маршала json перед отправкой пользователю
type Balance struct {
	Current 	float64 `json:"current"`
	Withdrawn 	float64 `json:"withdrawn"`
}