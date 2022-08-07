package pg

// Инициируем подключение к БД и возвращаем ссылку на объект

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
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
    							password VARCHAR (255) NOT NULL)`)
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

// UserExist проверяет наличие пользователя в БД по логину
func (p *pg) UserExist(ctx context.Context, login string) (bool, error) {
	var loginFromDB string
	err := p.db.QueryRowContext(ctx, `SELECT login FROM app_user WHERE login=$1 LIMIT 1`, login).Scan(&loginFromDB)
	// Записи нет в БД
	if fmt.Sprintf("%s", err) == "sql: no rows in result set" {
		return false, nil
	}
	// Обрабатываем прочие ошибки
	if err != nil {
		return false, fmt.Errorf("SQL Query Error: %s", err)
	}
	// Запись есть но логин не совпадает
	if login == loginFromDB {
		return true, nil
	}
	// Default
	return false, nil

}

// NewUser - создает нового пользователя и возвращает его id
func (p *pg) NewUser(ctx context.Context, login string, pwd string)  (int, error) {
	lastInsertId := 0
	err := p.db.QueryRow(`INSERT INTO app_user (login, password) VALUES ($1, $2) RETURNING id`, login, pwd).Scan(&lastInsertId)
	if err != nil {
		return 0, fmt.Errorf("new user insert error: %s", err)
	}
	return lastInsertId, nil
}