package pg

/*
	Конструктор объектов
	Создает
		- соединение с БД
		- HTTP клиент и запусает его по тикеру
		- Прочие взаиможействия с источниками данных
*/

import (
	"log"
)

// DBConfig добавляя поля в структуру можно настраивать приложение для подключения к конкретной СУБД.
// 			в нашем случает существует только Postgres по условию задачи.
type DBConfig struct {
	PGConnStr string
}

// NewDB возвращет ссылку на подключение к БД, инициируем схему.
func NewDB(conf DBConfig, logger *log.Logger) DBInterface {
	if len(conf.PGConnStr) != 0 {
		db := New(conf.PGConnStr)
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