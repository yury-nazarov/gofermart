package pg

import (
	"context"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
)

// DBInterface методы работы с релиационными БД
type DBInterface interface {
	Ping() bool
	// UserExist проверяет наличие пользователя в БД
	UserExist(ctx context.Context, login string) (bool, error)
	NewUser(ctx context.Context, login string, pwd string) (int, error)
	// UserIsValid Проверяет на сколько валидны креды пользователя и вообще существует ли он
	UserIsValid(ctx context.Context, login string, hashPwd string) (userID int, err error)
	// GetOrderByNumber Вернет заказ по его номеру
	GetOrderByNumber(ctx context.Context, orderNum string) (order models.OrderDB, err error)
	// AddOrder добавит новый номер заказа
	AddOrder(ctx context.Context, orderNumber string, userID int) error
	// AddAccrual добавляет запись в таблицу accrual
	AddAccrual(ctx context.Context, userID int) error
	// ListOrders Получить спосок заказов пользователя
	ListOrders(ctx context.Context, userID int) (orderList []models.OrderDB, err error)
	// GetOrders получает все заказы со статусом NEW, PROCESSING
	GetOrders() ([]string, error)
	// OrderStatusUpdate обновить статус заказа
	OrderStatusUpdate(ctx context.Context, orderNum string, status string) error
	//GetAccrual получить текущие значения таблицы: app_user.accrual_current, app_user.accrual_total
	GetAccrual(ctx context.Context, userID int) (models.UserDB, error)
	// UpdateAccrual - обновить значения таблицы: accrual.current_point, accrual.total_point
	UpdateAccrual(ctx context.Context, currentPoint float64, totalPoint float64, userID int) error
	// UpdateAccrualTransaction - списание баллов в рамках транзации и обновление / добавление во все нужные таблицы
	UpdateAccrualTransaction(ctx context.Context, orderNum string, userID int, sum float64) error
	// UpdateOrderAccrual - обновляет значения для app_order.accrual
	UpdateOrderAccrual(ctx context.Context, accrual float64, orderNumber string) error
	// GetOrderByUserID проверяем налицие заказа для конкретного пользователя
	GetOrderByUserID(ctx context.Context, orderNum string, userID int) (string, error)
	// AddToWithdrawList - добавляет новую запись в журнал
	AddToWithdrawList(ctx context.Context, orderNum string, sum float64, userID int) error
	// GetWithdrawList вернет список всех списаний для пользователя
	GetWithdrawList(ctx context.Context, userID int) ([]models.WithdrawDB, error)
}
