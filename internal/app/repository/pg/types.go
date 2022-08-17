package pg

import "context"

// OrderDB структура модели заказа в БД.
//		   Используется для описания таблицы в БД
type OrderDB struct {
	ID         int		`json:"-"`
	UserID     int		`json:"-"`
	Number     string   `json:"number"` 		// Номер заказа
	Status     string  	`json:"status"`			// Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    float64 	`json:"accrual"`		// Сколько начислено баллов этому заказу
	UploadedAt string  	`json:"uploaded_at"`	// Дата загрузки в формате RFC3339
}

// DBInterface методы работы с релиационными БД
type DBInterface interface {
	Ping() bool
	GetToken(ctx context.Context, token string) (bool, error)
	// UserExist проверяет наличие пользователя в БД
	UserExist(ctx context.Context, login string) (bool, error)
	NewUser(ctx context.Context, login string, pwd string) (int, error)
	// UserIsValid Проверяет на сколько валидны креды пользователя и вообще существует ли он
	UserIsValid(ctx context.Context, login string, hashPwd string) (userID int, err error)
	// GetOrderByNumber Вернет заказ по его номеру
	GetOrderByNumber(ctx context.Context, orderNum string) (order OrderDB, err error)
	// AddOrder добавит новый номер заказа
	AddOrder(ctx context.Context, orderNumber string, userID int) error
	// AddAccrual добавляет запись в таблицу accrual
	AddAccrual(ctx context.Context, userID int) error
	// ListOrders Получить спосок заказов пользователя
	ListOrders(ctx context.Context, userID int) (orderList []OrderDB, err error)
	// GetOrders получает все заказы со статусом NEW, PROCESSING
	GetOrders() ([]string, error)
	// OrderStatusUpdate обновить статус заказа
	OrderStatusUpdate(ctx context.Context, orderNum string, status string) error
	//GetAccrual получить текущие значения таблицы: app_user.accrual_current, app_user.accrual_total
	GetAccrual(ctx context.Context, userID int) (currentPoint float64, totalPoint float64, err error)
	// UpdateAccrual - обновить значения таблицы: accrual.current_point, accrual.total_point
	UpdateAccrual(ctx context.Context, currentPoint float64, totalPoint float64, userID int) error
	// UpdateOrderAccrual - обновляет значения для app_order.accrual
	UpdateOrderAccrual(ctx context.Context, accrual float64, orderNumber string) error
}
