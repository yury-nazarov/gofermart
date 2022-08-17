package withdraw

import (
	"context"
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
