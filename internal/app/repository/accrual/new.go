package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
)

// NewAccrual Создает объект для отправки запросов в систему рассчета баллов
func NewAccrual(accrualAddress string, db pg.DBInterface, logger *log.Logger) *accrualClientStruct {
	return &accrualClientStruct{
		accrualAddress: accrualAddress,
		db:             db,
		logger:         logger,
		ctx:            context.Background(),
	}
}

// Init запускает бесконечный цикл
func (a *accrualClientStruct) Init() {
	// 	В бесконечном цикле раз в полсекунды получаем из БД все заказы со статусом NEW, PROCESSING
	//	уточняем их состояние в accrual
	for {
		// Получаем все со статусом NEW, PROCESSING из БД
		orderList := a.getDataFromDB()
		for _, orderNum := range orderList {

			// Выполняем запрос в систему рассчета баллов
			order, err := a.getOrderByID(orderNum)
			if err != nil {
				continue
			}

			// Обновляем результат в БД
			if len(order.Number) != 0 {
				err := a.updateAccrual(order)
				if err != nil {
					a.logger.Printf("updateAccrual have error execute: %s", err)
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// getOrderByID получает данные из accrual системы
func (a *accrualClientStruct) getOrderByID(orderNum string) (models.OrderFromAccrualSystem, error) {
	order := models.OrderFromAccrualSystem{}

	endpoint := fmt.Sprintf("%s/api/orders/%s", a.accrualAddress, orderNum)
	resp, err := http.Get(endpoint)
	if err != nil {
		err = fmt.Errorf("can't connection to accrual server: %s", a.accrualAddress)
		return order, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		payload, err := io.ReadAll(resp.Body)
		if err != nil {
			//a.logger.Printf("can't read http body: %s", err)
			return order, err
		}

		err = json.Unmarshal(payload, &order)
		if err != nil {
			a.logger.Printf("HTTP Client unmarshal err %s", err)
		}
		return order, nil
	}
	return order, fmt.Errorf("order not fount")
}

// getDataFromDB - получает из БД заказы со стратусом NEW и PROCESSING
func (a *accrualClientStruct) getDataFromDB() []string {
	orders, err := a.db.GetOrders()
	if err != nil {
		a.logger.Printf("get orders from DB: NEW, PROCESSING is error: %s", err)
	}
	return orders
}

// updateAccrual
//			1. Получаем детальные данные о заказе по orderNum
//			2. Проверяем если статус заказа вернувшийся от системы рассчета боллов: PROCESSED -
//			   значит обработка завершена и можно начислить баллы пользователю и записать
//			   сколько баллов получено за конкретный заказ.
//			3. Проверяем если статус заказа вернувшийся от системы рассчета боллов: INVALID, PROCESSING, REGISTERED
//			   Если статус иной, то обновляем статус для заказа и идем дальше.
func (a *accrualClientStruct) updateAccrual(order models.OrderFromAccrualSystem) error {
	a.logger.Printf("begin processed order: %s", order.Number)
	//// Получить данные о текущем заказе и пользователе из БД
	orderDB, err := a.db.GetOrderByNumber(a.ctx, order.Number)
	if err != nil {
		errMsg := fmt.Errorf("can't get order by orderNumber: %s", err)
		a.logger.Println(errMsg)
		return errMsg
	}
	a.logger.Printf("success get order info from table 'app_order': num: '%s', status: '%s', accrual: '%f', userID: '%d'", orderDB.Number, orderDB.Status, orderDB.Accrual, orderDB.UserID)

	// Если статус PROCESSED - значит обработка завершена и получен результат с новым accrual
	// 						   можно обновить все необходимые данные в БД
	if order.Status == "PROCESSED" {
		// Посчитать новые значения для accrual.current_point, accrual.total_point
		user, err := a.db.GetAccrual(a.ctx, orderDB.UserID)
		if err != nil {
			errMsg := fmt.Errorf("can't get accrual on userID: '%d', err: %s", orderDB.UserID, err)
			a.logger.Println(errMsg)
			return errMsg
		}
		// При успешном получениее данных из accrual начисляем баллы
		user.AccrualCurrent += order.Accrual
		user.AccrualTotal += order.Accrual

		// Обновляем данные в таблице accrual.current_point, accrual.total_point для userID
		err = a.db.UpdateAccrual(a.ctx, user)
		if err != nil {
			errMsg := fmt.Errorf("can't update accrual for userID: '%d', err: %s", orderDB.UserID, err)
			a.logger.Println(errMsg)
			return errMsg
		}
		a.logger.Printf("success update table 'accrual' for userID: '%d'", orderDB.UserID)

		// Обновить данные об в таблице app_order.accrual для orderNumber
		err = a.db.UpdateOrderAccrual(a.ctx, order)
		if err != nil {
			errMsg := fmt.Errorf("can't update 'app_order.accrual' for userID: '%d', err: %s", orderDB.UserID, err)
			a.logger.Println(errMsg)
			return errMsg
		}

	}
	// Обновляем статус если: INVALID, PROCESSING, REGISTERED
	err = a.db.OrderStatusUpdate(a.ctx, order)
	if err != nil {
		errMsg := fmt.Errorf("can't update 'app_order.status' %s for order: '%s', err: %s", order.Status, order.Number, err)
		a.logger.Println(errMsg)
		return errMsg
	}
	return nil
}
