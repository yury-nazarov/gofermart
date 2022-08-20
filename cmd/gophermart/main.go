package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yury-nazarov/gofermart/internal/app/handler"
	"github.com/yury-nazarov/gofermart/internal/app/repository/accrual"
	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"github.com/yury-nazarov/gofermart/internal/app/service/auth"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"github.com/yury-nazarov/gofermart/internal/app/service/withdraw"
	"github.com/yury-nazarov/gofermart/pkg/logger"
)

func main() {
	// Устанавливаем логгер
	logger := logger.NewLogger("gofermart")

	// Иницииреуем необходимые переменные для работы сервиса из аргументов или env
	serverAddress, accrualAddress, pgConfig := initParams(logger)

	// Инициируем БД и создаем соединение
	db, err := pg.NewDB(pgConfig, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Инициируем loginCache для проверки сессии пользователя
	loginSession := cache.NewLoginCache()

	// Регистрация и авторизация пользователя
	user := auth.NewAuth(db, loginSession, logger)

	// Запускаем горутины в бусконечном цикле которые будут периодически опрашивать accrualServer и обновлять значение в БД
	accrualClient := accrual.NewAccrual(accrualAddress, db, logger)
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
	logger.Fatal(http.ListenAndServe(serverAddress, router))
}

// initRunArgs - Иницииреует необходимые переменные из аргументов или env
func initParams(logger *log.Logger) (string, string, string) {
	// Парсим аргументы командной строки
	serverAddressFlag := flag.String("a", "", "set server address, by example: 127.0.0.1:8080")
	accrualSystemFlag := flag.String("r", "", "set accrual server address, by example: 127.0.0.1")
	dbFlag := flag.String("d", "", "set database URI for Postgres, by example: host=localhost port=5432 user=example password=123 dbname=example sslmode=disable connect_timeout=5")
	flag.Parse()

	// Получаем переменные окружения
	serverAddressEnv := os.Getenv("RUN_ADDRESS")
	accrualSystemEnv := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	dbEnv := os.Getenv("DATABASE_URI")

	// Устанавливаем конфигурационные параметры по приоритету:
	// 		1. Флаги;
	// 		2. Переменные окружения;
	// Если данных не достаточно, завершаем работу приложения
	serverAddress, err := serverConfigInit(*serverAddressFlag, serverAddressEnv)
	if err != nil {
		logger.Fatalf("serverAddress: %s", err)
	}
	accrualAddress, err := serverConfigInit(*accrualSystemFlag, accrualSystemEnv)
	if err != nil {
		logger.Fatalf("accrualAddress: %s", err)
	}
	pgConfig, err := serverConfigInit(*dbFlag, dbEnv)
	if err != nil {
		logger.Fatalf("pgConfig: %s", err)
	}
	return serverAddress, accrualAddress, pgConfig
}

// serverConfigInit - возвращает приоритетное значение из переданых аргументов
func serverConfigInit(flag string, env string) (string, error) {
	if len(flag) != 0 {
		return flag, nil
	}
	if len(env) != 0 {
		return env, nil
	}
	return "", fmt.Errorf("you should set flag or env vars for config service")
}
