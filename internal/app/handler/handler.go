package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yury-nazarov/gofermart/pkg/tools"
	"io"
	"log"
	"net/http"

	"github.com/yury-nazarov/gofermart/internal/app/repository/accrual"
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"github.com/yury-nazarov/gofermart/internal/app/service/auth"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"github.com/yury-nazarov/gofermart/internal/app/service/withdraw"
)

type Controller struct {
	user         auth.UserInterface
	loginSession cache.UserSessionInterface
	order        processing.OrderInterface
	balance      withdraw.BalanceInterface
	accrual      accrual.AccrualInterface
	logger       *log.Logger
}

// New Создаем новый контроллер, через который будем управлять хендлерами
func New(user auth.UserInterface, loginSession cache.UserSessionInterface, order processing.OrderInterface,
	balance withdraw.BalanceInterface, accrual accrual.AccrualInterface, logger *log.Logger) *Controller {
	c := &Controller{
		user:         user,
		loginSession: loginSession,
		order:        order,
		balance:      balance,
		accrual:      accrual,
		logger:       logger,
	}
	return c
}

// Register регистрация пользователя
// 			200 — пользователь успешно аутентифицирован;
//			400 — неверный формат запроса;
//			401 — неверная пара логин/пароль;
//			500 — внутренняя ошибка сервера.
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	// Читаем присланые данные
	user := auth.User{}
	err400 := JSONError400(r, &user, c.logger)
	if err400 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Передаем в cлой бизнес логики
	// Регистрируем пользователя
	token, err409, err500 := c.user.SignUp(r.Context(), user.Login, user.Password)
	if err409 != nil {
		c.logger.Printf("can't sing up userLogin: %s, err: %s", user.Login, err409)
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err500 != nil {
		c.logger.Print(err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Отправялем ответ клиенту, записав токен в заголовок
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

// Login аутентификация пользователя
// 		200 — пользователь успешно аутентифицирован;
//		400 — неверный формат запроса;
//		401 — неверная пара логин/пароль или пользователь не существует;
//		500 — внутренняя ошибка сервера.
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	// Читаем присланые данные
	user := auth.User{}
	err400 := JSONError400(r, &user, c.logger)
	if err400 != nil {
		c.logger.Printf("JSON parsing error for login: %s, err: %s", user.Login, err400)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Передаем в лой бизнес логики
	token, err401, err500 := c.user.SignIn(r.Context(), user.Login, user.Password)
	if err401 != nil {
		c.logger.Printf("can't sign in login: %s, err: %s", user.Login, err400)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err500 != nil {
		c.logger.Print(err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправялем ответ клиенту, записав токен в заголовок
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

// AddOrders загрузка пользователем номера заказа для расчёта
// 			200 — номер заказа уже был загружен этим пользователем;
//			202 — новый номер заказа принят в обработку;
//			400 — неверный формат запроса;
//			401 — пользователь не аутентифицирован;
//			409 — номер заказа уже был загружен другим пользователем;
//			422 — неверный формат номера заказа;
//			500 — внутренняя ошибка сервера.
func (c *Controller) AddOrders(w http.ResponseWriter, r *http.Request) {
	// Читаем присланые данные из HTTP приводим к строке номер заказа
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Printf("HTTP Body parsing error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order := string(bodyData)


	// Получаем пользователя по токену
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil { // Ошибка подключения к кешу
		c.logger.Printf("can't connection to cache userID: %d. err: %s", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//// пользователь не авторизован (если по каким то причинам кеш с сессиями протух)
	//if userID == 0 {
	//	c.logger.Printf("can't authorisation userID: %d", userID)
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	var err409 tools.Error409
	var err422 tools.Error422
	var err500 tools.Error500
	// Пробуем добавить заказ
	ok200, ok202, err := c.order.Add(r.Context(), order, userID)
	// номер заказа уже был загружен этим пользователем;
	if ok200 {
		c.logger.Printf("Order %s for userID %d is exist", order, userID)
		w.WriteHeader(http.StatusOK)
		return
	}
	// новый номер заказа принят в обработку;
	if ok202 {
		c.logger.Printf("Order %s for userID %d accepted and will be processing", order, userID)
		w.WriteHeader(http.StatusAccepted)
		return
	}
	// номер заказа уже был загружен другим пользователем;
	//if err409 != nil {
	if errors.As(err, &err409) {
		c.logger.Printf("order exist, err409: %s", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	// неверный формат номера заказа
	//if err422 != nil {
	if errors.As(err, &err422) {
		c.logger.Printf("incorrect order format, err4: %s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// внутренняя ошибка сервера
	//if err500 != nil {
	if errors.As(err, &err500) {
		c.logger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


// GetOrders получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
//			200 — успешная обработка запроса.
// 			[
//      		{
//          		"number": "9278923470",
//          		"status": "PROCESSED",
//          		"accrual": 500,
//         		 	"uploaded_at": "2020-12-10T15:15:45+03:00"
//      		},
//				...
//			]
// 			204 — нет данных для ответа.
//			401 — пользователь не авторизован.
//			500 — внутренняя ошибка сервера.
func (c *Controller) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя по токену
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil { // Ошибка подключения к кешу
		c.logger.Printf("can't connection to cache userID: %d, err: %s", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: НУжна эта проверка?
	if userID == 0 { // пользователь не авторизован (если по каким то причинам кеш с сессиями протух)
		c.logger.Printf("can't authorisation userID: %d", userID)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Пробуем получить заказы пользователя
	var err204 tools.Error204
	var err500 tools.Error500
	//orders, err204, err500 := c.order.List(r.Context(), userID)
	orders, err := c.order.List(r.Context(), userID)
	//if err204 != nil {
	if errors.As(err, &err204) {
		c.logger.Printf("Order list is empty.  userID: %d, err: %s", userID, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if errors.As(err, &err500) {
		c.logger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Сериализуем JSON и отдаем пользователю
	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		c.logger.Printf("can't json marshal. err: %s", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersJSON)
	if err != nil {
		c.logger.Printf("can't write JSON to client. err: %s", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// GetBalance получение текущего баланса счёта баллов лояльности пользователя
// 			200 — успешная обработка запроса.
// 			401 — пользователь не авторизован.
//			500 — внутренняя ошибка сервера.
func (c *Controller) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя по токену
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil {
		errMsg := fmt.Errorf("can't connection to cache of userID: %d session: err %s", userID, err)
		c.logger.Print(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получаем текущий баланс
	balance, err := c.balance.CurrentBalance(r.Context(), userID)
	if err != nil {
		errMsg := fmt.Errorf("can't get current ballance. err: %s", err)
		c.logger.Print(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Преобразуем в JSON
	jsonBalance, err := json.Marshal(balance)
	if err != nil {
		errMsg := fmt.Errorf("can't convert ballance struct to JSON. error: %s", err)
		c.logger.Print(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBalance)
	if err != nil {
		errMsg := fmt.Errorf("can't send json to client. error: %s", err)
		c.logger.Print(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Withdraw запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
// 			200 — успешная обработка запроса;
//			401 — пользователь не авторизован;
//			402 — на счету недостаточно средств;
//			422 — неверный номер заказа;
//			500 — внутренняя ошибка сервера.
//          ------------------------------------
// 			POST /api/user/balance/withdraw HTTP/1.1
//			Content-Type: application/json
//			{
//    			"order": "2377225624",
//    			"sum": 751
//			}
func (c *Controller) Withdraw(w http.ResponseWriter, r *http.Request) {
	// Читаем присланые данные
	withdraw := pg.WithdrawDB{}
	err := JSONError400(r, &withdraw, c.logger)
	if err != nil {
		c.logger.Printf("can't json read. err: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получить из заголовка token и преобразовать его в userID
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil {
		errMsg := fmt.Errorf("can't connection to cache of userID: %d session: err %s", userID, err)
		c.logger.Print(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Выводим средства со счета пользователя: app_user.current - sum
	var err402 tools.Error402
	var err422 tools.Error422
	var err500 tools.Error500
	//err402, err422, err500 := c.balance.WithdrawBalance(r.Context(), userID, withdraw.Order, withdraw.Sum)
	err = c.balance.WithdrawBalance(r.Context(), userID, withdraw.Order, withdraw.Sum)
	//if err402 != nil {
	if errors.As(err, &err402) {
		c.logger.Printf("can't calculate withdraw balance: err %s", err)
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	if errors.As(err, &err422){
		c.logger.Printf("order number wrong: err %s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if errors.As(err, &err500) {
		c.logger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Withdrawals получение информации о выводе средств с накопительного счёта пользователем
// 					200 — успешная обработка запроса.
// 							200 OK HTTP/1.1
//  						Content-Type: application/json
//  						...
//  						[
//      						{
//          						"order": "2377225624",
//          						"sum": 500,
//          						"processed_at": "2020-12-09T16:09:57+03:00"
//      						}
//  						]
// 					204 — нет ни одного списания.
//					401 — пользователь не авторизован.
//					500 — внутренняя ошибка сервера.
func (c *Controller) Withdrawals(w http.ResponseWriter, r *http.Request) {
	// Получить из заголовка token и преобразовать его в userID
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil {
		errMsg := fmt.Errorf("can't connection to cache of userID: %d session: err %s", userID, err)
		c.logger.Print(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var err204 tools.Error204
	var err500 tools.Error500

	withdrawList, err := c.balance.Withdrawals(r.Context(), userID)
	if errors.As(err, &err204) {
		c.logger.Printf("withdraw list for userID: %d is empty", userID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if errors.As(err, &err500) {
		c.logger.Printf("can't connection to cache of user session: err %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Сериализуем JSON и отдаем пользователю
	withdrawListJSON, err := json.Marshal(withdrawList)
	if err != nil {
		c.logger.Printf("can't json marshal. err: %s\n", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("withdrawList", withdrawList)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(withdrawListJSON)
	if err != nil {
		c.logger.Printf("can't send json to client. error: %s", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
