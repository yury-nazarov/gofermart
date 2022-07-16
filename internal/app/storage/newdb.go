package storage

import (
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/storage/pg"
	"github.com/yury-nazarov/gofermart/internal/app/storage/repository"
)

// DBConfig добавляя поля в структуру можно настраивать приложение для подключения к конкретной СУБД.
// 			в нашем случает существует только Postgres по условию задачи.
type DBConfig struct {
	PGConnStr	string
}

// New возвращет ссылку на подключение к БД, инициируем схему.
func New(conf DBConfig, logger *log.Logger) repository.Repository {
	if len(conf.PGConnStr) != 0 {
		db := pg.New(conf.PGConnStr)
		// Проверяем соединение с БД
		if !db.Ping() {
			logger.Fatal("DB not connected. Ping fail.")
		}
		// Применяем схему
		if err := db.SchemeInit(); err != nil {
			logger.Fatalf("Postgres DB init is err: %s", err)
		}
		logger.Println("DB Postgres is connecting")
		return db
	}
	logger.Fatal("DB not selected")
	return nil
}