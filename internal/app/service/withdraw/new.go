package withdraw

import (
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"log"
)

type balanceStruct struct {
	db     pg.DBInterface
	logger *log.Logger
}

func NewBalance(db pg.DBInterface, logger *log.Logger) BalanceInterface {
	return balanceStruct{
		db:     db,
		logger: logger,
	}
}
