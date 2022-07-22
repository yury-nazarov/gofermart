package pg

// Инициируем подключение к БД и возвращаем ссылку на объект

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"

)

// pg объект через который происходит подключение к БД
type pg struct {
	db *sql.DB
}

// New Иницирует подключение к Postgres
func New(connStr string) *pg {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("sql.Open is err: %s", err)
	}
	dbConnect := &pg{
		db: db,
	}
	return dbConnect
}

// SchemeInit создает схему БД если ее еще нет
func (p *pg) SchemeInit() error {
	return nil
}

// Ping - Проверка соединения с БД
func (p *pg) Ping() bool {
	if err := p.db.Ping(); err != nil {
		log.Printf("Ping fail:, %s", err)
		return false
	}
	return true
}

func (p *pg) GetToken(ctx context.Context, token string) (bool, error) {
	return true, nil
}