package handler

import (
	"log"
	"net/http"

	"github.com/yury-nazarov/gofermart/internal/app/repository"
)

type Controller struct {
	db     		repository.DBInterface
	accrual 	repository.AccrualInterface
	logger 		*log.Logger
	// as accrualService
}

// New объект через который получаем доступ к основным ручкам обслуживающим энедпоинты
func New(db repository.DBInterface, accrual repository.AccrualInterface, logger *log.Logger) *Controller {
	c := &Controller{
		db:     	db,
		accrual: 	accrual,
		logger: 	logger,
	}
	return c
}

// Register регистрация пользователя
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {

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
