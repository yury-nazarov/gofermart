package accrual

import (
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
)

// New Создает клиент для отправки запросов в систему рассчета баллов
func New(accrualAddress string, db pg.DBInterface, logger *log.Logger) AccrualInterface {
	return nil
}
