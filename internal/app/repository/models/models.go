package models

// Модели для описания структуры данных во внешних системах и взаимодествия с ними:
// - Accrual - система рассчета баллов

// OrderFromAccrualSystem данные которые мы получаем из системы рассчета баллов
// 						  используем для анмаршала HTTP Request отвеа от системы рассчета баллов
// 						  OrderFromAccrualSystem и OrderDB отличаются json полем Number "order" vs "number"
//						  это ломает ответы пользователю.
type OrderFromAccrualSystem struct {
	Number  string  `json:"order"`   // Номер заказа
	Status  string  `json:"status"`  // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual float64 `json:"accrual"` // Сколько начислено баллов этому заказу
}

// Модели представления данных в СУБД и для формирования ответов пользователям при обращении в HTTP ручку

// OrderDB структура данных для
//			таблицы: app_order
//			передачи в методы
//			маршала JSON в HTTP GET для ручки handler.GetOrders
type OrderDB struct {
	ID         int     `json:"-"`
	UserID     int     `json:"-"`
	Number     string  `json:"number"`      // Номер заказа
	Status     string  `json:"status"`      // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    float64 `json:"accrual"`     // Сколько начислено баллов этому заказу
	UploadedAt string  `json:"uploaded_at"` // Дата загрузки в формате RFC3339
}

// WithdrawDB структура данных для
//			таблицы: withdraw_list
//			анмаршала JSON из HTTP POST для ручки handler.Withdraw
//			маршала JSON в HTTP GET для ручки handler.Withdrawals
// 			передачи в методы
type WithdrawDB struct {
	ID          int     `json:"-"`
	UserID      int     `json:"-"`
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

// UserDB структура данных для таблицы: app_user
// 			Используем для передачи в методы
type UserDB struct {
	ID             int     `json:"-"`
	Login          string  `json:"login,omitempty"`
	Password       string  `json:"password,omitempty"`
	Token          string  `json:"-"`
	AccrualCurrent float64 `json:"accrual_current"`
	AccrualTotal   float64 `json:"accrual_total"`
}
