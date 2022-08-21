package withdraw

import (
	"context"
	"fmt"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"github.com/yury-nazarov/gofermart/pkg/tools"
)

type balanceStruct struct {
	db     pg.DBInterface
	logger *log.Logger
}

func NewBalance(db pg.DBInterface, logger *log.Logger) *balanceStruct {
	return &balanceStruct{
		db:     db,
		logger: logger,
	}
}

func (b *balanceStruct) CurrentBalance(ctx context.Context, userID int) (Balance, error) {
	var balance Balance
	// Делаем запрос в app_user, получаем: app_user.accrual_current app_user.accrual_total
	current, total, err := b.db.GetAccrual(ctx, userID)
	if err != nil {
		return balance, err
	}
	// Считаем сетрики которые хочет увидеть пользователь
	balance.Current = current
	balance.Withdrawn = total - current

	return balance, nil
}

func (b *balanceStruct) WithdrawBalance(ctx context.Context, userID int, orderNum string, sum float64) error {
	// Проверить номер заказа Луном
	err := processing.CorrectOrderNumber(orderNum)
	if err != nil {
		errMgg := fmt.Sprintf("incorrect order number '%s'. err: %s", orderNum, err.Error())
		return tools.NewError422(errMgg)
	}

	err = b.db.UpdateAccrualTransaction(ctx, orderNum, userID, sum)
	if err != nil {
		// Транзитом прокидываем 402 и 500 на верх
		return err
	}

	return nil
}

func (b *balanceStruct) Withdrawals(ctx context.Context, userID int) (WithdrawList []models.WithdrawDB, err error) {
	// Получить данные из таблицы withdraw_list
	RawWithdrawList, err := b.db.GetWithdrawList(ctx, userID)
	if err != nil {
		return nil, tools.NewError500(err.Error())
	}
	if len(RawWithdrawList) == 0 {
		return nil, tools.NewError204("'withdraw_list' is empty'")
	}
	// Преобразовать дату в RFC3339
	for _, v := range RawWithdrawList {
		dataRFC3339, err := tools.ToRFC3339(v.ProcessedAt, "Europe/Moscow")
		if err != nil {
			b.logger.Printf("can't convert datatime. err: %s", err)
		}
		withdraw := models.WithdrawDB{}
		withdraw.Order = v.Order
		withdraw.Sum = v.Sum
		withdraw.ProcessedAt = dataRFC3339

		WithdrawList = append(WithdrawList, withdraw)

	}
	return WithdrawList, nil
}
