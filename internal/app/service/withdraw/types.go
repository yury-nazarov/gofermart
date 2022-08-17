package withdraw

import "context"

/*
	Типы для работы со слоем бизнес логики: вывод средств и запрос текущего баланса
*/

type BalanceInterface interface {
	// CurrentBalance возвращает текущий баланс
	CurrentBalance(ctx context.Context, userID int) (Balance, error)
	// WithdrawBalance выводит средства со счета польователя
	WithdrawBalance(ctx context.Context, userID int, order string, sum float64) (err402 error, err422 error, err500 error)
	// WriteToWithdrawList Заносим информацию о списании в журнал
	WriteToWithdrawList(ctx context.Context, orderNum string, sum float64) error
	// ReadFromWithdrawList Получаем все записи из журнала
	ReadFromWithdrawList(ctx context.Context, orderNum string) (Withdraw, error)
}

// Balance Для маршала json перед отправкой пользователю
type Balance struct {
	Current 	float64 `json:"current"`
	Withdrawn 	float64 `json:"withdrawn"`
}

// Withdraw структура для анмаршала JSON из HTTP Request
type Withdraw struct {
	 Order 			string 		`json:"order"`
	 Sum 			float64 	`json:"sum"`
	 ProcessedAt 	string 		`json:"processed_at,omitempty"`
}