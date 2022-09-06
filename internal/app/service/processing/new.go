package processing

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"github.com/yury-nazarov/gofermart/pkg/tools"

	"github.com/theplant/luhn"
)

func NewOrder(db pg.DBInterface, logger *log.Logger) orderStruct {
	return orderStruct{
		db:     db,
		logger: logger,
	}
}

// Add - добавляет новый заказ
func (o orderStruct) Add(ctx context.Context, newOrder models.OrderDB) (ok200, ok202 bool, err error) {
	// err422 - Проверяем корректен ли номер заказа
	// если номер заказа некорректный - отвечаем со статусом 422
	err = CorrectOrderNumber(newOrder.Number)
	if err != nil {
		// 422
		errMsg := fmt.Sprintf("orderNum: '%s', correctOrderNumber is wrong: %s", newOrder.Number, err)
		return false, false, tools.NewError422(errMsg)
	}

	// Проверяем наличие номера заказа в БД, а так же соответствие userID
	orderDB, err := o.db.GetOrderByNumber(ctx, newOrder.Number)

	// Если такого заказа нет - его можно создать
	// ok202 - заказ принят в обработку
	if err != nil {
		// Не получилось, подумать. go vet:  second argument to errors.As should not be *error
		//if errors.As(err, &sql.ErrNoRows) {
		err = o.db.AddOrder(ctx, newOrder)
		if err != nil {
			// err500
			errMsg := fmt.Sprintf("add order. err: %s", err)
			return false, false, tools.NewError500(errMsg)
		}
		// ok202 - заказ принят в обработку
		return false, true, nil
	}

	// err409 - пользователь уже добавил этот заказ
	if orderDB.Number == newOrder.Number && orderDB.UserID == newOrder.UserID {
		return true, false, nil
	}

	// err409 - Заказ создан другим пользователем
	if orderDB.Number == newOrder.Number && orderDB.UserID != newOrder.UserID {
		// 409
		return false, false, tools.NewError409("order exist for other user")
	}

	return false, false, tools.NewError500("create order something wrong")
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
func (o orderStruct) List(ctx context.Context, userID int) (orders []models.OrderDB, err error) {
	// Делаем запрос в БД
	orders, err = o.db.ListOrders(ctx, userID)
	// 	err500 - ошибка выполнения запроса
	if err != nil {
		errMsg := fmt.Sprintf("listOrder SQL Error: %s", err)
		return orders, tools.NewError500(errMsg)
	}
	// 	err204 - список пуст
	if len(orders) == 0 {
		return orders, tools.NewError204("empty order list")
	}
	orders = o.orderConvertData(orders)
	// ok200
	return orders, nil
}

// orderConvertData - конвертирует определенные поля заказа в нужный формат
func (o orderStruct) orderConvertData(orderList []models.OrderDB) (clearOrderList []models.OrderDB) {

	log.Println("Order list", orderList)
	for _, order := range orderList {

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

		// Добавляем результат в новый слайс
		clearOrderList = append(clearOrderList, order)
	}
	return clearOrderList
}
