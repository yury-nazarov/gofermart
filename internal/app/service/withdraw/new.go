package withdraw

import (
	"context"
	"fmt"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
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


func (b *balanceStruct) CurrentBalance(ctx context.Context, userID int) (Balance, error){
	var balance Balance
	// Делаем запрос в app_user, получаем: app_user.accrual_current app_user.accrual_total
	//current, total, err := b.db.GetUserBalance(ctx, userID)
	current, total, err := b.db.GetAccrual(ctx, userID)
	if err != nil {
		return balance, err
	}
	// Считаем сетрики которые хочет увидеть пользователь
	balance.Current = current
	balance.Withdrawn = total - current

	return balance, nil
}

func (b *balanceStruct) WithdrawBalance(ctx context.Context, userID int, orderNum string, sum float64) (err402 error, err422 error, err500 error) {
	// TODO: Возможно надо как то по другому проверять на  422
	//// err422:  Проверить наличие заказа у пользователя если его нет err402 - заказ не найден
	//_, err := b.db.GetOrderByUserID(ctx, orderNum, userID)
	//if err != nil {
	//	return nil, fmt.Errorf("order %s not found", orderNum), nil
	//}
	// Проверить номер заказа Луном
	err := processing.CorrectOrderNumber(orderNum)
	if err != nil {
		return nil, fmt.Errorf("order %s not found", orderNum), nil
	}

	// Получить текущее значение app_user.accrual_current
	accrualCurrent, accrualTotal, err := b.db.GetAccrual(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("can't get accrual. err: %s", err)
	}
	// err402: Не достаточно средств
	b.logger.Printf("DEBUG: compare %f < %f", accrualCurrent, sum)
	if accrualCurrent < sum {
		return fmt.Errorf("not enough points"), nil, nil
	}


	// Посчитать app_user.accrual_current - sum
	newAccrualCurrent := accrualCurrent - sum
	// записать в app_user.accrual_current
	err = b.db.UpdateAccrual(ctx, newAccrualCurrent, accrualTotal, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("can't update accrual. err: %s", err)
	}

	// записать в withdraw_list
	err = b.db.AddToWithdrawList(ctx, orderNum, sum)
	if err != nil {
		return nil, nil, fmt.Errorf("can't insert to withdraw_list. err: %s", err)
	}
	return nil, nil, nil
}


func (b *balanceStruct) WriteToWithdrawList(ctx context.Context, orderNum string, sum float64) error {
	return nil
}

func (b *balanceStruct) ReadFromWithdrawList(ctx context.Context, orderNum string) (Withdraw, error) {
	var w Withdraw
	return w, nil
}