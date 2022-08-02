package service

import (
	"github.com/yury-nazarov/gofermart/internal/app/repository"
	"log"
)

type balanceStruct struct {
	db repository.DBInterface
	logger *log.Logger
}


func NewBalance(db repository.DBInterface, logger *log.Logger) BalanceInterface {
	return balanceStruct{
		db: db,
		logger: logger,
	}
}