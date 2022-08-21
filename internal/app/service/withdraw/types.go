package withdraw

import (
	"context"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
)

/*
	Типы для работы со слоем бизнес логики: вывод средств и запрос текущего баланса
*/

type BalanceInterface interface {
	// CurrentBalance возвращает текущий баланс
	CurrentBalance(ctx context.Context, userID int) (Balance, error)
	// WithdrawBalance выводит средства со счета польователя
	//WithdrawBalance(ctx context.Context, userID int, order string, sum float64) (err error)
	WithdrawBalance(ctx context.Context, withdrawal models.WithdrawDB) error
	// Withdrawals - возвращает список списаний для пользователя
	Withdrawals(ctc context.Context, userID int) (WithdrawList []models.WithdrawDB, err error)
}

// Balance Для маршала json перед отправкой пользователю
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
