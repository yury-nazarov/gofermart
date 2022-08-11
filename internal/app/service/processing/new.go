package processing

import (
	"context"
	"fmt"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"log"

	"github.com/theplant/luhn" //	алгоритм Луна для проверки корректности номера
)

func NewOrder(db pg.DBInterface, logger *log.Logger) orderStruct {
	return orderStruct{
		db:     db,
		logger: logger,
	}
}

// Add - добавляет новый заказ
func (o orderStruct) Add(ctx context.Context, orderNum int, userID int) (ok200, ok202 bool, err409, err422, err500 error) {
	// err422 - Проверяем корректен ли номер заказа
	// если номер заказа некорректный - отвечаем со статусом 422
	if !luhn.Valid(orderNum) {
		return false, false, nil, fmt.Errorf("wrong order number format"), nil
	}

	// Проверяем наличие номера заказа в БД, а так же соответствие userID

	order, err := o.db.GetOrderByNumber(ctx, orderNum)

	// Если произошла ошибка, то такого заказа нет и его можно создать
	// ok202 - заказ принят в обработку
	if err != nil {
		err500 = o.db.AddOrder(ctx, orderNum, userID)
		if err500 != nil {
			// err500
			return false, false, nil, nil, fmt.Errorf("add order errer %s", err500)
		}
		// ok202 - заказ принят в обработку
		return false, true, nil, nil, nil
	}


	// err409 - пользователь уже добавил этот заказ
	if order.Number == orderNum && order.UserID == userID {
		return true, false, nil, nil, nil
	}

	// err409 - Заказ создан другим пользователем
	if order.Number == orderNum && order.UserID != userID {
		return false, false, fmt.Errorf("order exist for other user"), nil, nil
	}

	log.Printf("create order somfing wrong")
	return false, false, nil, nil, fmt.Errorf("create order somfing wrong")
}

// List - список всех заказов пользователя
func (o orderStruct) List(ctx context.Context, userID int) (orders []pg.OrderDB, err204, err500 error) {
	orders = []pg.OrderDB{}
	// Делаем запрос в БД
	// 	204 - список пуст
	// 	500 - ошибка выполнения запроса
	//  orders - добавляем в список объекты для дальнейшей сериализации
	return orders, nil, nil
}
