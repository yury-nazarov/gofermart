package pg

// Инициируем подключение к БД и возвращаем ссылку на объект

import (
	"context"
	"database/sql"
	"fmt"
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
	// Контекст для инициализации БД
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Таблица Users - содержит логин пользователя и хеш пароля
	_, err := p.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS app_user (
    							id serial PRIMARY KEY,
    							login VARCHAR (255) NOT NULL,
    							password VARCHAR (255) NOT NULL,
    							token VARCHAR (255) NOT NULL)`)
	if err != nil {
		return fmt.Errorf("create table `user`: %w", err)
	}
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

// GetUser вернет данные пользователя по логину
func (p *pg) GetUser(ctx context.Context, login string) (bool, error) {
	return true, nil
}

// NewUser - создает нового пользователя
func (p *pg) NewUser(ctx context.Context, login string, pwd string) (string, error) {
	return "", nil
}