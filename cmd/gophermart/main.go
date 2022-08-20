package main

import (
	"net/http"

	"github.com/yury-nazarov/gofermart/internal/app/handler"
	"github.com/yury-nazarov/gofermart/internal/app/repository/accrual"
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/repository/config"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"github.com/yury-nazarov/gofermart/internal/app/service/auth"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"github.com/yury-nazarov/gofermart/internal/app/service/withdraw"
	"github.com/yury-nazarov/gofermart/pkg/logger"
)

func main() {
	// Устанавливаем логгер
	logger := logger.NewLogger("gofermart")

	// Инициируем конфиг: аргументы cli > env
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// Инициируем БД и создаем соединение
	db, err := pg.NewDB(cfg.DB, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Инициируем loginCache для проверки сессии пользователя
	loginSession := cache.NewLoginCache()

	// Регистрация и авторизация пользователя
	user := auth.NewAuth(db, loginSession, logger)

	// Запускаем горутины в бусконечном цикле которые будут периодически опрашивать accrualServer и обновлять значение в БД
	accrualClient := accrual.NewAccrual(cfg.AccrualSystemAddress, db, logger)
	go accrualClient.Init()

	// Бизнес логика работы с заказами
	order := processing.NewOrder(db, logger)

	// Бизнес логика работы с балансом пользователя
	balance := withdraw.NewBalance(db, logger)

	// Инициируем объект для доступа к хендлерам
	c := handler.New(user, loginSession, order, balance, accrualClient, logger)

	// инициируем роутер
	router := handler.NewRouter(c, user, logger)

	// Запускаем сервер
	logger.Fatal(http.ListenAndServe(cfg.RunAddress, router))
}
