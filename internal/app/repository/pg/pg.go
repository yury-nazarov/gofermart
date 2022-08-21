package pg

// Инициируем подключение к БД и возвращаем ссылку на объект

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"

	"github.com/yury-nazarov/gofermart/internal/app/repository/models"
	"github.com/yury-nazarov/gofermart/pkg/tools"

	"github.com/pressly/goose"
)

// pg объект через который происходит подключение к БД
type pg struct {
	db *sql.DB
	logger *log.Logger
}

// New Иницирует подключение к Postgres
func New(connStr string, logger *log.Logger) (*pg, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open is err: %s", err)
	}
	dbConnect := &pg{
		db: db,
		logger: logger,
	}
	return dbConnect, nil
}

func (p *pg) Migration() error {
	err := goose.Up(p.db, "./internal/migrations")
	if err != nil {
		return err
	}
	return nil
}


// Ping - Проверка соединения с БД
func (p *pg) Ping() bool {
	if err := p.db.Ping(); err != nil {
		log.Printf("Ping fail:, %s", err)
		return false
	}
	p.logger.Printf("success ping db check")
	return true
}

// GetUser вернет данные пользователя по логину

// UserExist проверяет наличие пользователя в БД по логину
func (p *pg) UserExist(ctx context.Context, login string) (bool, error) {
	var loginFromDB string
	err := p.db.QueryRowContext(ctx, `SELECT login FROM app_user WHERE login=$1 LIMIT 1`, login).Scan(&loginFromDB)
	// Записи нет в БД
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	// Обрабатываем прочие ошибки
	if err != nil {
		return false, fmt.Errorf("SQL Query Error: %s", err)
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
func (p *pg) GetOrderByNumber(ctx context.Context, orderNum string) (o models.OrderDB, err error) {
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

// AddOrder добавит новый номер заказа
func (p *pg) AddOrder(ctx context.Context, orderNum string, userID int) (err500 error) {
	_, err500 = p.db.ExecContext(ctx, `INSERT INTO app_order (number, user_id, status, accrual) 
                                             VALUES ($1, $2, $3, $4)`, orderNum, userID, "NEW", 0)
	if err500 != nil {
		return err500
	}
	return nil
}

// AddAccrual добавляет запись в таблицу accrual
func (p *pg) AddAccrual(ctx context.Context, userID int) error {
	_, err500 := p.db.ExecContext(ctx, `INSERT INTO accrual (current_point, total_point, user_id) 
                                              VALUES (0, 0, $1)`, userID)
	if err500 != nil {
		return err500
	}
	return nil
}

// ListOrders Получить спосок заказов пользователя
func (p *pg) ListOrders(ctx context.Context, userID int) (orderList []models.OrderDB, err error) {
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

	o := models.OrderDB{}
	for rows.Next() {
		log.Println("Upload", &o.UploadedAt)
		if err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			log.Println("error read string for order list")
		}
		orderList = append(orderList, o)
	}
	return orderList, nil
}

// GetOrders получает все заказы со статусом NEW, PROCESSING
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
	for rows.Next() {
		var order string
		if err = rows.Scan(&order); err != nil {
			log.Println(err)
		} else {
			orderList = append(orderList, order)
		}

	}
	return orderList, nil
}

// OrderStatusUpdate обновить статус заказа
func (p *pg) OrderStatusUpdate(ctx context.Context, orderNum string, status string) error {
	_, err500 := p.db.ExecContext(ctx, `UPDATE app_order SET status=$1 WHERE number=$2`, status, orderNum)
	if err500 != nil {
		return err500
	}
	return nil
}

//GetAccrual получить текущие значения таблицы: app_user.accrual_current, app_user.accrual_total
func (p *pg) GetAccrual(ctx context.Context, userID int) (currentPoint float64, totalPoint float64, err error) {

	err = p.db.QueryRowContext(ctx, `SELECT accrual_current, accrual_total FROM app_user
                                           WHERE id=$1 LIMIT 1`, userID).Scan(&currentPoint, &totalPoint)
	if err != nil {
		return 0, 0, err
	}
	return currentPoint, totalPoint, nil
}

// UpdateAccrual - обновить значения таблицы: accrual.current_point, accrual.total_point
func (p *pg) UpdateAccrual(ctx context.Context, currentPoint float64, totalPoint float64, userID int) error {
	_, err := p.db.ExecContext(ctx, `UPDATE app_user
                                           SET accrual_current=$1, accrual_total=$2
                                           WHERE id=$3`, currentPoint, totalPoint, userID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrderAccrual - обновляет значения для app_order.accrual
func (p *pg) UpdateOrderAccrual(ctx context.Context, accrual float64, orderNumber string) error {
	_, err := p.db.ExecContext(ctx, `UPDATE app_order SET accrual=$1
                                           WHERE number=$2`, accrual, orderNumber)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAccrualTransaction - обновить значения таблиц: app_user, app_order
func (p *pg) UpdateAccrualTransaction(ctx context.Context, orderNum string, userID int, sum float64) error {
	// Открываем транзакцию
	tx, err := p.db.Begin()
	if err != nil {
		errMsg := fmt.Sprintf("can't open transaction. err: %s", err)
		return tools.NewError500(errMsg)
	}
	defer tx.Rollback()

	// Получаем данные из БД о текущем балансе пользователя
	var accrualCurrent float64
	err = tx.QueryRowContext(ctx, `SELECT accrual_current FROM app_user WHERE id=$1 LIMIT 1`, userID).Scan(&accrualCurrent)
	if err != nil {
		errMsg := fmt.Sprintf("transaction select user accrual has err: %s", err)
		return tools.NewError500(errMsg)
	}

	// err402: Не достаточно средств
	if accrualCurrent < sum {
		return tools.NewError402("not enough points")
	}

	// Посчитать app_user.accrual_current - sum
	newAccrualCurrent := accrualCurrent - sum

	// Готовим стейтмент для апдейта app_user
	updateAccrual, err := tx.PrepareContext(ctx, "UPDATE app_user SET accrual_current=$1 WHERE id=$2")
	if err != nil {
		errMsg := fmt.Sprintf("transaction statment updateAccrual has err: %s", err)
		return tools.NewError500(errMsg)
	}
	defer updateAccrual.Close()

	// Готовим стейтмент для апдейта withdraw_list
	updateWithdrawList, err := tx.PrepareContext(ctx, "INSERT INTO withdraw_list (order_num, sum_points, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		errMsg := fmt.Sprintf("transaction statment updateWithdrawList has err: %s", err)
		return tools.NewError500(errMsg)
	}

	// Выполянем
	_, err = updateAccrual.ExecContext(ctx, newAccrualCurrent, userID)
	if err != nil {
		errMsg := fmt.Sprintf("transaction execute updateAccrual has err: %s", err)
		return tools.NewError500(errMsg)
	}

	// Выполянем
	_, err = updateWithdrawList.ExecContext(ctx, orderNum, sum, userID)
	if err != nil {
		errMsg := fmt.Sprintf("transaction execute updateOrderAccrual has err: %s", err)
		return tools.NewError500(errMsg)
	}

	// Применяем
	return tx.Commit()
}

// GetOrderByUserID проверяем налицие заказа для конкретного пользователя
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

// AddToWithdrawList - добавляет новую запись в журнал
func (p *pg) AddToWithdrawList(ctx context.Context, orderNum string, sumPoints float64, userID int) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO withdraw_list (order_num, sum_points, user_id) 
                                           VALUES ($1, $2, $3)`, orderNum, sumPoints, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetWithdrawList вернет список всех списаний для пользователя
func (p *pg) GetWithdrawList(ctx context.Context, userID int) (withdrawList []models.WithdrawDB, err error) {
	rows, err := p.db.QueryContext(ctx, `SELECT order_num, sum_points, processed_at
                                               FROM withdraw_list
                                               WHERE user_id=$1
                                               ORDER BY processed_at`, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Printf("defer rows.Close() error: %s", err)
		}
		err = rows.Err()
		if err != nil {
			log.Printf("defer rows.Err() error: %s", err)
		}
	}()
	withdraw := models.WithdrawDB{}
	for rows.Next() {
		err = rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			log.Printf("can't read string for withdrow_list: %s", err)
		}
		withdrawList = append(withdrawList, withdraw)
	}
	return withdrawList, nil
}
