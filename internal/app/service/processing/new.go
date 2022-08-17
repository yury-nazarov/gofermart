package processing

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"

	"github.com/theplant/luhn"
)

func NewOrder(db pg.DBInterface, logger *log.Logger) orderStruct {
	return orderStruct{
		db:     db,
		logger: logger,
	}
}

// Add - добавляет новый заказ
func (o orderStruct) Add(ctx context.Context, orderNum string, userID int) (ok200, ok202 bool, err409, err422, err500 error) {
	// err422 - Проверяем корректен ли номер заказа
	// если номер заказа некорректный - отвечаем со статусом 422
	err422 = CorrectOrderNumber(orderNum)
	if err422 != nil {
		return false, false, nil, fmt.Errorf("orderNum: '%s', correctOrderNumber is wrong: %s", orderNum, err422), nil
	}


	// Проверяем наличие номера заказа в БД, а так же соответствие userID
	order, err := o.db.GetOrderByNumber(ctx, orderNum)

	// Если произошла ошибка, то такого заказа нет и его можно создать
	// ok202 - заказ принят в обработку
	if err != nil {
		err500 = o.db.AddOrder(ctx, orderNum, userID)
		if err500 != nil {
			// err500
			return false, false, nil, nil, fmt.Errorf("add order error %s", err500)
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

// CorrectOrderNumber - проверяет корректность номера заказа
// 						Длина больше 0 и коректен по луну
func CorrectOrderNumber(orderNum string) error {
	if len(orderNum) == 0 {
		return fmt.Errorf("order num is empty")
	}
	luhnCheck, err := strconv.Atoi(orderNum)
	if err != nil {
		return fmt.Errorf("strconv.Atoi err, %s", err)
	}
	if !luhn.Valid(luhnCheck) {
		return fmt.Errorf("wrong luhn order number format")
	}
	return nil
}

// List - список всех заказов пользователя
func (o orderStruct) List(ctx context.Context, userID int) (orders []pg.OrderDB, err204, err500 error) {
	// Делаем запрос в БД
	orders, err := o.db.ListOrders(ctx, userID)
	// 	err500 - ошибка выполнения запроса
	if err != nil {
		return nil, nil, fmt.Errorf("listOrder SQL Error: %s", err)
	}
	// 	err204 - список пуст
	if len(orders) == 0 {
		return nil, fmt.Errorf("empty order list"), nil
	}
	orders = o.orderConvertData(orders)
	// ok200
	return orders, nil, nil
}

// orderConvertData - конвертирует определенные поля заказа в нужный формат
func (o orderStruct) orderConvertData(orderList []pg.OrderDB) (clearOrderList []pg.OrderDB) {

	log.Println("Order list", orderList)
	for _, order := range orderList {

		// TODO: По хорошему нужно из конфига, но этого нет в ТЗ :(
		// Устанавливаем тайм зону
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			o.logger.Printf("load location error %s", err)
		}

		// Из строки получаем объект Time в нужном формате и локейшене
		newOrderTime, err := time.ParseInLocation(time.RFC3339, order.UploadedAt, loc)
		if err != nil {
			o.logger.Printf("error conver time %s", err)
		}

		// Форматируем в: "2020-12-10T15:15:45+03:00"
		order.UploadedAt = newOrderTime.In(loc).Format("2006-01-2T15:04:05Z07:00")

		// TODO: Пока не переложил в новый слайс, новое значение не возвращалось. Разобратся в чем проблема!
		// Добавляем результат в новый слайс
		clearOrderList = append(clearOrderList, order)
	}
	return clearOrderList
}