package processing

import (
	"context"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
)

/*
	Типы для работы со слоем бизнес логики: Обработка заказа
*/

// orderStruct - структура возвращаемая конструктором
type orderStruct struct {
	db     pg.DBInterface
	logger *log.Logger
}

// OrderInterface - интерфейс логики работы с заказом на базе структуры: orderStruct
type OrderInterface interface {
	// Add Добавить заказ
	//Add(ctx context.Context, order string, userID int) (ok200, ok202 bool, err409, err422, err500 error)
	Add(ctx context.Context, order string, userID int) (ok200, ok202 bool, err error)
	// List Получить список заказов
	List(ctx context.Context, userID int) (orders []pg.OrderDB, err error)
}
