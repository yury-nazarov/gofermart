package repository

///*
//	Описывает интерфейсы работы со слоем
//	Работа с источниками данных
//*/
//
//import (
//	"context"
//	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
//)
//
//// DBInterface методы работы с релиационными БД
//type DBInterface interface {
//	Ping() bool
//	GetToken(ctx context.Context, token string) (bool, error)
//	// UserExist проверяет наличие пользователя в БД
//	UserExist(ctx context.Context, login string) (bool, error)
//	NewUser(ctx context.Context, login string, pwd string) (int, error)
//	// UserIsValid Проверяет на сколько валидны креды пользователя и вообще существует ли он
//	UserIsValid(ctx context.Context, login string, hashPwd string) (userID int, err error)
//	// GetOrderByNumber Вернет заказ по его номеру
//	GetOrderByNumber(ctx context.Context, orderNum int) (order pg.OrderDB, err error)
//	// AddOrder добавит новый номер заказа
//	AddOrder(ctx context.Context, orderNumber int, userID int) error
//}
//
//// AccrualInterface методы работы с системой рассчета баллов
//type AccrualInterface interface {
//
//}
