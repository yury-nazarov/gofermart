package сlient

import (
	"context"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
)

// AccrualInterface методы работы с системой рассчета баллов
type AccrualInterface interface {
	Init()
	getOrder(orderNum string) (orderNumber string, orderStatus string, orderAccrual float64)
	getDataFromDB() []string
	updateAccrual(orderNumber string, orderStatus string, orderAccrual float64) error
}

// accrualClientStruct создает объект для работы с системой рассчета баллов
type accrualClientStruct struct {
	accrualAddress 	string
	db 				pg.DBInterface
	logger 			*log.Logger
	ctx 			context.Context
}

// AccrualOrder для преобразования из JSON ответа accrual сервиса
type AccrualOrder struct {
	Number 	string 	`json:"order"`
	Status 	string 	`json:"status"`
	Accrual float64 `json:"accrual"`
}


