package handler

import (
	"github.com/yury-nazarov/gofermart/internal/app/repository/auth"
	"github.com/yury-nazarov/gofermart/internal/app/service"
	"log"
	"net/http"

	"github.com/yury-nazarov/gofermart/internal/app/repository"
)


type Controller struct {
	//db     		repository.DBInterface
	user		auth.UserInterface
	order 		service.OrderInterface
	balance 	service.BalanceInterface
	accrual 	repository.AccrualInterface
	logger 		*log.Logger
	// as accrualService
}


func New(user auth.UserInterface, order service.OrderInterface, balance service.BalanceInterface, accrual repository.AccrualInterface, logger *log.Logger) *Controller {
	c := &Controller{
		user: 		user,
		order: 		order,
		balance: 	balance,
		accrual: 	accrual,
		logger: 	logger,
	}
	return c
}

// Register регистрация пользователя
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	//// Читаем присланые данные
	//bodyData , err := io.ReadAll(r.Body)
	//if err != nil || len(bodyData) == 0 {
	//	c.logger.Printf("the HTTP Body parsing error: %w", err)
	//	w.WriteHeader(http.StatusBadRequest)
	//}
	//
	//// Unmarshal JSON
	//var user auth.User
	//if err = json.Unmarshal(bodyData, &user); err != nil {
	//	c.logger.Printf("unmarshal json error: %w", err)
	//	w.WriteHeader(http.StatusBadRequest)
	//}
	//
	//// Проверяем наличие пользователя в БД
	//ok, err := auth.IsUserExist(r.Context(), user.Login)
	//if err != nil {
	//	c.logger.Printf("executing isUserExist() error: %w", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//}
	//// Пользователь уже существует
	//if ok {
	//	c.logger.Printf("%s user exist", user.Login)
	//	w.WriteHeader(http.StatusConflict)
	//}
	//// Добавляем пользователя
	//token, err := auth.NewUserToken(r.Context(), user.Login, user.Password)
	//if err != nil {
	//	c.logger.Printf("executing newUserToken() error: %w", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//}
	//
	//// Отправялем ответ клиенту, записав токен в заголовок
	//w.Header().Set("Authorization", token)
	//
	//if err != nil {
	//	c.logger.Printf("write HTTP answer error: %w", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//}
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
