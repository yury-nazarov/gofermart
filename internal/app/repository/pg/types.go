package pg

import "context"

// OrderDB структура модели заказа в БД.
//		   Используется для описания таблицы в БД
type OrderDB struct {
	Id         int
	UserID     int
	Number     int     // Номер заказа
	Status     string  // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    float64 // Сколько начислено баллов этому заказу
	UploadedAt string  // Дата загрузки в формате RFC3339
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
	GetOrderByNumber(ctx context.Context, orderNum int) (order OrderDB, err error)
	// AddOrder добавит новый номер заказа
	AddOrder(ctx context.Context, orderNumber int, userID int) error
}
