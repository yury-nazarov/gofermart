package handler

import (
	"log"
	"net/http"

	"github.com/yury-nazarov/gofermart/internal/app/repository"
	"github.com/yury-nazarov/gofermart/internal/app/repository/auth"
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/service"
)


type Controller struct {
	//db     		repository.DBInterface
	user			auth.UserInterface
	loginSession	cache.UserSessionInterface
	order 			service.OrderInterface
	balance 		service.BalanceInterface
	accrual 		repository.AccrualInterface
	logger 			*log.Logger
	// as accrualService
}


func New(user auth.UserInterface, loginSession cache.UserSessionInterface, order service.OrderInterface, balance service.BalanceInterface, accrual repository.AccrualInterface, logger *log.Logger) *Controller {
	c := &Controller{
		user: 			user,
		loginSession: 	loginSession,
		order: 			order,
		balance: 		balance,
		accrual: 		accrual,
		logger: 		logger,
	}
	return c
}

// TODO: 04/08/22 Реализуем ручку для регимстрации пользователя: Register

// Register регистрация пользователя
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	// Читаем присланые данные
	user := auth.User{}
	err400 := c.JSONError400(r, &user)
	if err400 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
	return
}

// Login аутентификация пользователя
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {

}

// AddOrders загрузка пользователем номера заказа для расчёта
func (c *Controller) AddOrders(w http.ResponseWriter, r *http.Request) {

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
