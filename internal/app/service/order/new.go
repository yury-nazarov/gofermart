package order

import (
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository"
)

type orderStruct struct {
	db repository.DBInterface
	logger *log.Logger
}


func NewOrder(db repository.DBInterface, logger *log.Logger) OrderInterface {
	return orderStruct{
		db: db,
		logger: logger,
	}
}