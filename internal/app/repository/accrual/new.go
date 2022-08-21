package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
	"io"
	"log"
	"net/http"
	"time"

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
		a.logger.Printf("HTTP Client: get Orders from DB: %s", orderList)
		for _, order := range orderList {
			a.logger.Printf("HTTP Client: try to connect accrual server: %s for get info about order: %s", a.accrualAddress, order)

			// Выполняем запрос в систему рассчета баллов
			//orderNum, status, accrual, err := a.getOrder(order)
			order, err := a.getOrderByID(order)
			if err != nil {
				a.logger.Printf("can't connect to accrual system. err: %s", err)
				continue
			}
			a.logger.Printf("orderNum: %s, status: %s, accrual: %f", order.Number, order.Status, order.Accrual)

			// Обновляем результат в БД
			if len(order.Number) != 0 {
				a.logger.Printf("success get data from accrual system: orderNum: %s, status: %s, accrual: %f\n",  order.Number, order.Status, order.Accrual)
				err := a.updateAccrual( order.Number, order.Status, order.Accrual)
				if err != nil {
					a.logger.Printf("updateAccrual have error execute: %s", err)
				}
			}
		}
		a.logger.Println("accrual.Init()----------------------------------")
		time.Sleep(500 * time.Millisecond)
	}
}

// getOrder получает данные из accrual системы
//func (a *accrualClientStruct) getOrder(orderNum string) (string, string, float64, error) {
func (a *accrualClientStruct) getOrderByID(orderNum string) (models.OrderDB, error) {
	order := models.OrderDB{}

	endpoint := fmt.Sprintf("%s/api/orders/%s", a.accrualAddress, orderNum)
	a.logger.Printf("HTTP Client: HTTP GET to endpoint: %s", endpoint)
	resp, err := http.Get(endpoint)
	if err != nil {
		err = fmt.Errorf("can't connection to accrual server: %s", a.accrualAddress)
		a.logger.Print(err)
		return order, err
	}
	defer resp.Body.Close()

	a.logger.Printf("HTTP Client: response status code: %d", resp.StatusCode)
	//a.logger.Printf("HTTP Client: response body: %d", resp.Body)
	if resp.StatusCode == 200 {
		payload, err := io.ReadAll(resp.Body)
		if err != nil {
			a.logger.Printf("can't read http body: %s", err)
			return order, err
		}
		//var order AccrualOrder

		err = json.Unmarshal(payload, &order)
		if err != nil {
			a.logger.Printf("HTTP Client unmarshal err %s", err)
		}
		//return order.Number, order.Status, order.Accrual, nil
		return order, nil
	}
	return order, nil
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
func (a *accrualClientStruct) updateAccrual(orderNum string, status string, accrual float64) error {
	a.logger.Printf("begin processed order: %s", orderNum)
	//// Получить данные о текущем заказе и пользователе из БД
	order, err := a.db.GetOrderByNumber(a.ctx, orderNum)
	if err != nil {
		errMsg := fmt.Errorf("can't get order by orderNumber: %s", err)
		a.logger.Println(errMsg)
		return errMsg
	}
	a.logger.Printf("success get order info from table 'app_order': num: '%s', status: '%s', accrual: '%f', userID: '%d'", order.Number, order.Status, order.Accrual, order.UserID)

	// Если статус PROCESSED - значит обработка завершена и получен результат с новым accrual
	// 						   можно обновить все необходимые данные в БД
	if status == "PROCESSED" {
		// Посчитать новые значения для accrual.current_point, accrual.total_point
		currentPoint, totalPoint, err := a.db.GetAccrual(a.ctx, order.UserID)
		if err != nil {
			errMsg := fmt.Errorf("can't get accrual on userID: '%d', err: %s", order.UserID, err)
			a.logger.Println(errMsg)
			return errMsg
		}
		// При успешном получениее данных из accrual начисляем баллы
		currentPoint += accrual
		totalPoint += accrual

		// Обновляем данные в таблице accrual.current_point, accrual.total_point для userID
		err = a.db.UpdateAccrual(a.ctx, currentPoint, totalPoint, order.UserID)
		if err != nil {
			errMsg := fmt.Errorf("can't update accrual for userID: '%d', err: %s", order.UserID, err)
			a.logger.Println(errMsg)
			return errMsg
		}
		a.logger.Printf("success update table 'accrual' for userID: '%d'", order.UserID)

		// Обновить данные об в таблице app_order.accrual для orderNumber
		err = a.db.UpdateOrderAccrual(a.ctx, accrual, order.Number)
		if err != nil {
			errMsg := fmt.Errorf("can't update 'app_order.accrual' for userID: '%d', err: %s", order.UserID, err)
			a.logger.Println(errMsg)
			return errMsg
		}

	}
	// Обновляем статус если: INVALID, PROCESSING, REGISTERED
	err = a.db.OrderStatusUpdate(a.ctx, orderNum, status)
	if err != nil {
		errMsg := fmt.Errorf("can't update 'app_order.status' %s for order: '%s', err: %s", status, orderNum, err)
		a.logger.Println(errMsg)
		return errMsg
	}
	return nil
}
