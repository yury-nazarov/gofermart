package processing

import (
	"context"

	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
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
	// Add Добавить заказ
	Add(ctx context.Context, newOrder models.OrderDB) (ok200, ok202 bool, err error)
	// List Получить список заказов
	List(ctx context.Context, userID int) (orders []models.OrderDB, err error)
}
