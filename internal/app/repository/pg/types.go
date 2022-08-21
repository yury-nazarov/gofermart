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
	// NewUser регистрация нового пользователя
	NewUser(ctx context.Context, user models.UserDB) (int, error)
	// UserIsValid Проверяет на сколько валидны креды пользователя и вообще существует ли он
	UserIsValid(ctx context.Context, login string, hashPwd string) (userID int, err error)
	// GetOrderByNumber Вернет заказ по его номеру
	GetOrderByNumber(ctx context.Context, orderNum string) (order models.OrderDB, err error)
	// AddOrder добавит новый номер заказа
	AddOrder(ctx context.Context, order models.OrderDB) error
	// AddAccrual добавляет запись в таблицу accrual
	AddAccrual(ctx context.Context, userID int) error
	// ListOrders Получить спосок заказов пользователя
	ListOrders(ctx context.Context, userID int) (orderList []models.OrderDB, err error)
	// GetOrders получает все заказы со статусом NEW, PROCESSING
	GetOrders() ([]string, error)
	// OrderStatusUpdate обновить статус заказа
	OrderStatusUpdate(ctx context.Context, order models.OrderFromAccrualSystem) error
	//GetAccrual получить текущие значения таблицы: app_user.accrual_current, app_user.accrual_total
	GetAccrual(ctx context.Context, userID int) (models.UserDB, error)
	// UpdateAccrual - обновить значения таблицы: accrual.current_point, accrual.total_point
	UpdateAccrual(ctx context.Context, user models.UserDB) error
	// UpdateAccrualTransaction - списание баллов в рамках транзации из таблицы 'app_user'
	// 							  добавление записи о списании в 'withdraw_list'
	UpdateAccrualTransaction(ctx context.Context, withdrawal models.WithdrawDB) error
	// UpdateOrderAccrual - обновляет значения для app_order.accrual
	UpdateOrderAccrual(ctx context.Context, order models.OrderFromAccrualSystem) error
	// GetWithdrawList вернет список всех списаний для пользователя
	GetWithdrawList(ctx context.Context, userID int) ([]models.WithdrawDB, error)
}
