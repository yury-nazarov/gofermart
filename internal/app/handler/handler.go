package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/yury-nazarov/gofermart/internal/app/repository/accrual"
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/service/auth"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"github.com/yury-nazarov/gofermart/internal/app/service/withdraw"
)

type Controller struct {
	//db     		repository.DBInterface
	user         	auth.UserInterface
	loginSession 	cache.UserSessionInterface
	order        	processing.OrderInterface
	balance withdraw.BalanceInterface
	accrual accrual.AccrualInterface
	logger  *log.Logger
	// as accrualService
}

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
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err500 != nil {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Передаем в лой бизнес логики
	token, err401, err500 := c.user.SignIn(r.Context(), user.Login, user.Password)
	if err401 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err500 != nil {
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
//			+ 400 — неверный формат запроса;
//			+ 401 — пользователь не аутентифицирован;
//			409 — номер заказа уже был загружен другим пользователем;
//			422 — неверный формат номера заказа;
//			500 — внутренняя ошибка сервера.
func (c *Controller) AddOrders(w http.ResponseWriter, r *http.Request) {
	// Читаем присланые данные из HTTP приводим к строке номер заказа
	bodyData, err400 := io.ReadAll(r.Body)
	if err400 != nil {
		c.logger.Printf("handlers/AddOrders, err400: HTTP Body parsing error: %s", err400)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order := string(bodyData)


	// Получаем пользователя по токену
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil { // Ошибка подключения к кешу
		c.logger.Printf("handlers/AddOrders, userID: %d can't connection to cache\n", userID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// пользователь не авторизован (если по каким то причинам кеш с сессиями протух)
	if userID == 0 {
		c.logger.Printf("handlers/AddOrders, userID: %d not authorisation\n", userID)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Пробуем добавить заказ
	ok200, ok202, err409, err422, err500 := c.order.Add(r.Context(), order, userID)
	// номер заказа уже был загружен этим пользователем;
	if ok200 {
		c.logger.Printf("handlers/AddOrders, ok200: Order %s for userID %d is exist\n", order, userID)
		w.WriteHeader(http.StatusOK)
		return
	}
	// новый номер заказа принят в обработку;
	if ok202 {
		c.logger.Printf("handlers/AddOrders, ok202: Order %s for userID %d accepted and will be processing\n", order, userID)
		w.WriteHeader(http.StatusAccepted)
		return
	}
	// номер заказа уже был загружен другим пользователем;
	if err409 != nil {
		c.logger.Printf("handlers/AddOrders, err409: %s\n", err409)
		w.WriteHeader(http.StatusConflict)
		return
	}
	c.logger.Printf("DEBUG orderNum: '%s'\n", order)
	// неверный формат номера заказа
	if err422 != nil {
		c.logger.Printf("handlers/AddOrders, err422: %s\n", err422)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// внутренняя ошибка сервера
	if err500 != nil {
		c.logger.Printf("handlers/AddOrders, err500: %s\n", err500)
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
		c.logger.Printf("handlers/GetOrders, userID: %d can't connection to cache\n", userID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.logger.Printf("handlers/GetOrders, got userID: %d, for token: %s", userID, token)

	// TODO: НУжна эта проверка?
	if userID == 0 { // пользователь не авторизован (если по каким то причинам кеш с сессиями протух)
		c.logger.Printf("handlers/GetOrders, userID: %d not authorisation\n", userID)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Пробуем получить заказы пользователя
	orders, err204, err500 := c.order.List(r.Context(), userID)
	if err204 != nil {
		c.logger.Printf("handlers/GetOrders,err204:  userID: %d. Order list is empty\n", userID)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err500 != nil {
		c.logger.Printf("handlers/GetOrders, from c.order.List got err500: %s\n", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Сериализуем JSON и отдаем пользователю
	ordersJSON, err500 := json.Marshal(orders)
	if err500 != nil {
		c.logger.Printf("handlers/GetOrders, json marshal return err500: %s\n", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err500 = w.Write(ordersJSON)
	if err500 != nil {
		c.logger.Printf("handlers/GetOrders, w.Write return err500: %s\n", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// GetBalance получение текущего баланса счёта баллов лояльности пользователя
// 			200 — успешная обработка запроса.
// 			+ 401 — пользователь не авторизован.
//			500 — внутренняя ошибка сервера.
func (c *Controller) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя по токену
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil {
		errMsg := fmt.Errorf("can't connection to cache of user session: err %s", err)
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
		errMsg := fmt.Errorf("can't send json bites to client. error: %s", err)
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
	withdraw := withdraw.Withdraw{}
	err400 := JSONError400(r, &withdraw, c.logger)
	if err400 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c.logger.Printf("success get withdraw info. orderNum: %s, sum: %f", withdraw.Order, withdraw.Sum)

	// Получить из заголовка token и преобразовать его в userID
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil {
		errMsg := fmt.Errorf("can't connection to cache of user session: err %s", err)
		c.logger.Print(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.logger.Printf("success get userID from user session: %d", userID)

	// Выводим средства со счета пользователя: app_user.current - sum
	err402, err422, err500 := c.balance.WithdrawBalance(r.Context(), userID, withdraw.Order, withdraw.Sum)
	if err402 != nil {
		c.logger.Printf("can't calculate withdraw balance: err %s", err402)
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	if err422 != nil {
		c.logger.Printf("can't calculate withdraw balance: err %s", err422)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err500 != nil {
		c.logger.Printf("can't calculate withdraw balance: err %s", err500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.logger.Printf("success calculate withdraw balance for userID: %d", userID)
	w.WriteHeader(http.StatusOK)
	return
}

// Withdrawals получение информации о выводе средств с накопительного счёта пользователем
func (c *Controller) Withdrawals(w http.ResponseWriter, r *http.Request) {

}
