package models

// Модели для описания структуры данных во внешних системах и взаимодествия с ними:
// - Accrual - система рассчета баллов

// OrderFromAccrualSystem данные которые мы получаем из системы рассчета баллов
// 						  используем для десериализации данных отвеа системы рассчета баллов
// 						  OrderFromAccrualSystem и OrderDB отличаются json полем Number "order" vs "number"
//						  это ломает ответы пользователю.
type OrderFromAccrualSystem struct {
	Number  string  `json:"order"`   // Номер заказа
	Status  string  `json:"status"`  // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual float64 `json:"accrual"` // Сколько начислено баллов этому заказу
}

// Модели представления данных в СУБД и для формирования ответов пользователям при обращении в HTTP ручку

// OrderDB структура данных для таблицы: app_order
//			Используем для передачи в методы
//			Сериализации ответа ручки: handler.GetOrders
type OrderDB struct {
	ID         int     `json:"-"`
	UserID     int     `json:"-"`
	Number     string  `json:"number"`      // Номер заказа
	Status     string  `json:"status"`      // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    float64 `json:"accrual"`     // Сколько начислено баллов этому заказу
	UploadedAt string  `json:"uploaded_at"` // Дата загрузки в формате RFC3339
}

// WithdrawDB структура данных для таблицы: withdraw_list для анмаршала JSON из HTTP Request
type WithdrawDB struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

type UserDB struct {
	ID             int     `json:"-"`
	Login          string  `json:"-"`
	AccrualCurrent float64 `json:"accrual_current"`
	AccrualTotal   float64 `json:"accrual_total"`
}
