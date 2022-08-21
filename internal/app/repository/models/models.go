package models

// OrderDB структура модели заказа в БД.
//		   Используется для описания таблицы в БД
type OrderDB struct {
	ID         int     `json:"-"`
	UserID     int     `json:"-"`
	//Number     string  `json:"number"`      // Номер заказа
	Number     string  `json:"order"`       // Номер заказа
	Status     string  `json:"status"`      // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    float64 `json:"accrual"`     // Сколько начислено баллов этому заказу
	UploadedAt string  `json:"uploaded_at"` // Дата загрузки в формате RFC3339
}

//// AccrualOrder для преобразования из JSON ответа accrual сервиса
//type AccrualOrder struct {
//	Number  string  `json:"order"`
//	Status  string  `json:"status"`
//	Accrual float64 `json:"accrual"`
//}


// WithdrawDB структура для анмаршала JSON из HTTP Request
type WithdrawDB struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}