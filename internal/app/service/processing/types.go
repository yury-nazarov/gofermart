package processing

import (
	"context"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
)

/*
	Типы для работы со слоем бизнес логики: Обработка заказа
*/

// orderStruct - структура возвращаемая конструктором
type orderStruct struct {
	db     pg.DBInterface
	logger *log.Logger
}

// OrderInterface - интерфейс логики работы с заказом на базе структуры: orderStruct
type OrderInterface interface {
	Add(ctx context.Context, order int, userID int) (ok200, ok202 bool, err409, err422, err500 error) // Добавить заказ
	List(ctx context.Context, userID int) (orders []pg.OrderDB, err204, err500 error)                 // Список заказов
}

// Order структура описывает заказ
//		Используется для
//			описания таблицы в БД
//		    сериализации/десерилизации в/из JSON
type Order struct {
	Id         int     `json:"id,-"`
	UserID     int     `json:"user_id,-"`
	Number     int  `json:"order_id"`              // Номер заказа
	Status     string  `json:"status,omitempty"`      // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    float64 `json:"accrual,omitempty"`     // Сколько начислено баллов этому заказу
	UploadedAt string  `json:"uploaded_at,omitempty"` // Дата загрузки в формате RFC3339
}
