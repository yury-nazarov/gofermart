package handler

import (
	"github.com/yury-nazarov/gofermart/internal/app/repository/accrual"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"log"
	"net/http"

	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/service/auth"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"github.com/yury-nazarov/gofermart/internal/app/service/withdraw"
)

type Controller struct {
	//db     		repository.DBInterface
	user         auth.UserInterface
	loginSession cache.UserSessionInterface
	order        processing.OrderInterface
	balance      withdraw.BalanceInterface
	accrual      accrual.AccrualInterface
	logger       *log.Logger
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
	// Читаем и валидируем присланые данные
	//order := processing.Order{}
	order := pg.OrderDB{}

	err400 := JSONError400(r, &order.Number, c.logger)
	if err400 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Получаем пользователя по токену
	token := r.Header.Get("Authorization")
	userID, err := c.loginSession.GetUserIDByToken(token)
	if err != nil { // Ошибка подключения к кешу
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userID == 0 { // пользователь не авторизован (если по каким то причинам кеш с сессиями протух)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Пробуем добавить заказ
	ok200, ok202, err409, err422, err500 := c.order.Add(r.Context(), order.Number, userID)
	if ok200 { // номер заказа уже был загружен этим пользователем;
		w.WriteHeader(http.StatusOK)
		return
	}
	if ok202 { // новый номер заказа принят в обработку;
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if err409 != nil { // номер заказа уже был загружен другим пользователем;
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err422 != nil { // неверный формат номера заказа;
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err500 != nil { // внутренняя ошибка сервера.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// GetOrders получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (c *Controller) GetOrders(w http.ResponseWriter, r *http.Request) {

}

// GetBalance получение текущего баланса счёта баллов лояльности пользователя
func (c *Controller) GetBalance(w http.ResponseWriter, r *http.Request) {

}

// Withdraw запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (c *Controller) Withdraw(w http.ResponseWriter, r *http.Request) {

}

// Withdrawals получение информации о выводе средств с накопительного счёта пользователем
func (c *Controller) Withdrawals(w http.ResponseWriter, r *http.Request) {

}
