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
	_, errUser := p.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS app_user (
											id serial 			PRIMARY KEY,
											login 				VARCHAR (255) NOT NULL,
											password 			VARCHAR (255) NOT NULL,
    										accrual_current 	FLOAT default 0,
    										accrual_total 		FLOAT default 0
											)`)
	if errUser != nil {
		return fmt.Errorf("create table `app_user`: %s", errUser)
	}

	_, errOrder := p.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS app_order (
    										id serial 	PRIMARY KEY,
    										number 		VARCHAR (255) NOT NULL,
    										user_id 	INT NOT NULL,
    										status 		VARCHAR (255) NOT NULL,
    										accrual 	FLOAT,
    										uploaded_at TIMESTAMP default NOW()
    										)`)
	if errOrder != nil {
		return fmt.Errorf("create table `app_order`: %s", errOrder)
	}


	_, errWithdrawList := p.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS withdraw_list (
    											id serial 		PRIMARY KEY,
    											order_num	 	VARCHAR (255) NOT NULL,
    											sum_points	 	FLOAT,
    											user_id 		INT NOT NULL,
                                   				processed_at 	TIMESTAMP default NOW()
												)`)
	if errWithdrawList != nil {
		return fmt.Errorf("create table `withdraw_list`: %w", errWithdrawList)
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
func (p *pg) NewUser(ctx context.Context, login string, hashPwd string) (int, error) {
	lastInsertID := 0
	err := p.db.QueryRow(`INSERT INTO app_user (login, password) VALUES ($1, $2) RETURNING id`, login, hashPwd).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("new user insert error: %s", err)
	}
	fmt.Printf("lastInsertID %d for user %s\n", lastInsertID, login)
	return lastInsertID, nil
}

// UserIsValid - Делает SQL в БД если по login и хеш пароля есть запись - значит пользователь существует и валиден.
func (p *pg) UserIsValid(ctx context.Context, login string, hashPwd string) (userID int, err error) {
	err = p.db.QueryRowContext(ctx, `SELECT app_user.id FROM app_user
											WHERE login=$1
											AND password=$2
											LIMIT 1`, login, hashPwd).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("user not found: %s", err)
	}
	return userID, nil
}

// GetOrderByNumber Вернет заказ по его номеру
func (p *pg) GetOrderByNumber(ctx context.Context, orderNum string) (o OrderDB, err error) {
	// Если записи нет то вернется {0 0 0  0 }
	row := p.db.QueryRowContext(ctx, `SELECT id, user_id, number, status, accrual, uploaded_at 
										   FROM app_order 
										   WHERE number=$1 LIMIT 1`, orderNum)
	err = row.Scan(&o.UserID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
	if err != nil {
		return o, fmt.Errorf("order not found: %s", err)
	}
	return o, nil
}

// AddOrder -
func (p *pg) AddOrder(ctx context.Context, orderNum string, userID int) (err500 error) {
	_, err500 = p.db.ExecContext(ctx, `INSERT INTO app_order (number, user_id, status, accrual) 
											 VALUES ($1, $2, $3, $4)`, orderNum, userID, "NEW", 0)
	if err500 != nil {
		return err500
	}
	return nil
}

func (p *pg) AddAccrual(ctx context.Context, userID int) error {
	_, err500 := p.db.ExecContext(ctx, `INSERT INTO accrual (current_point, total_point, user_id) 
											  VALUES (0, 0, $1)`, userID)
	if err500 != nil {
		return err500
	}
	return nil
}


// ListOrders -
func (p *pg) ListOrders(ctx context.Context, userID int) (orderList []OrderDB, err error) {
	rows, err := p.db.QueryContext(ctx, `SELECT number, status, accrual, uploaded_at 
				 	   							FROM app_order WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("defer rows.Close() error")
		}
		err = rows.Err()
		if err != nil {
			log.Println("defer rows.Err()  error")
		}
	}()

	o := OrderDB{}
	for rows.Next() {
		log.Println("Upload", &o.UploadedAt)
		if err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			log.Println("error read string for order list")
		}
		orderList = append(orderList, o)
	}
	return orderList, nil
}

// TODO: Errors
func (p *pg) GetOrders() ([]string, error) {
	var orderList []string
	rows, err := p.db.Query(`SELECT number FROM app_order WHERE status='NEW' OR status='PROCESSING'`)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
		}
	}()
	for rows.Next(){
		var order string
		if err = rows.Scan(&order); err != nil {
			log.Println(err)
		} else {
			orderList = append(orderList, order)
		}

	}
	return orderList, nil
}

func (p *pg) OrderStatusUpdate(ctx context.Context, orderNum string, status string) error {
	_, err500 := p.db.ExecContext(ctx, `UPDATE app_order SET status=$1 WHERE number=$2`, status, orderNum)
	if err500 != nil {
		return err500
	}
	return nil
}

func (p *pg) GetAccrual(ctx context.Context, userID int) (currentPoint float64, totalPoint float64, err error) {

	err = p.db.QueryRowContext(ctx, `SELECT accrual_current, accrual_total FROM app_user
											WHERE id=$1 LIMIT 1`, userID).Scan(&currentPoint, &totalPoint)
	if err != nil {
		return 0, 0,  err
	}
	return currentPoint, totalPoint, nil
}

func (p *pg) UpdateAccrual(ctx context.Context, currentPoint float64, totalPoint float64, userID int) error {
	_, err := p.db.ExecContext(ctx, `UPDATE app_user 
										   SET accrual_current=$1, accrual_total=$2 
										   WHERE id=$3`, currentPoint, totalPoint, userID)
	if err != nil {
		return err
	}
	return nil
}

func (p *pg) UpdateOrderAccrual(ctx context.Context, accrual float64, orderNumber string) error {
	_, err := p.db.ExecContext(ctx, `UPDATE app_order SET accrual=$1
										   WHERE number=$2`, accrual, orderNumber)
	if err != nil {
		return err
	}
	return nil
}

func (p *pg) GetOrderByUserID(ctx context.Context, orderNum string, userID int) (string, error) {
	var status string

	err := p.db.QueryRowContext(ctx, `SELECT status FROM app_order 
											WHERE number=$1 
											AND user_id=$2 LIMIT 1`, orderNum, userID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}


func (p *pg) AddToWithdrawList(ctx context.Context, orderNum string, sumPoints float64, userID int) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO withdraw_list (order_num, sum_points, user_id) 
											VALUES ($1, $2, $3)`, orderNum, sumPoints, userID)
	if err != nil {
		return err
	}
	return nil
}