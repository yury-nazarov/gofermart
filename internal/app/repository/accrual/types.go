package accrual

import (
	"context"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
)

// AccrualInterface методы работы с системой рассчета баллов
type AccrualInterface interface {
	Init()
	getOrderByID(orderNum string) (db models.OrderFromAccrualSystem, err error)
	getDataFromDB() []string
	updateAccrual(db models.OrderFromAccrualSystem) error
}

// accrualClientStruct создает объект для работы с системой рассчета баллов
type accrualClientStruct struct {
	accrualAddress string
	db             pg.DBInterface
	logger         *log.Logger
	ctx            context.Context
}
