package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/yury-nazarov/gofermart/internal/app/accrual"
	"github.com/yury-nazarov/gofermart/internal/app/handler"
	"github.com/yury-nazarov/gofermart/internal/app/storage"
)

func main() {
	// Устанавливаем логгер
	logger := NewLogger()

	// Иницииреуем необходимые переменные для работы сервиса из аргументов или env
	serverAddress, accrualAddress, pgConfig := initParams(logger)

	// Инициируем БД и создаем соединение
	db := storage.New(storage.DBConfig{PGConnStr: pgConfig}, logger)

	// Запускаем по тикеру горутины которые будут периодически опрашивать accrualServer и обновлять значение в БД
	accrual.Start(accrualAddress, logger)

	// Инициируем объект для доступа к хендлерам
	c := handler.New(db, logger)

	// инициируем роутер
	router := newRouter(c)

	// Запускаем сервер
	logger.Fatal(http.ListenAndServe(serverAddress, router))
}

// NewLogger - создает новый логер
func NewLogger() *log.Logger {
	return log.New(os.Stdout, `gofermart | `, log.LstdFlags|log.Lshortfile)
}

// newRouter - создает роутер. Внутри определяем все ручки
func newRouter(c *handler.Controller) http.Handler {
	// Инициируем Router
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API endpoints
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", c.Register)
			r.Post("/login", c.Login)
			r.Route("/orders", func(r chi.Router) {
				r.Post("/", c.AddOrders)
				r.Get("/", c.GetOrders)
			})
			r.Route("/balance", func(r chi.Router) {
				r.Get("/", c.GetBalance)
				r.Get("/withdrawals", c.Withdrawals)
				r.Post("/withdraw", c.Withdraw)
			})
		})
	})
	return r
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
