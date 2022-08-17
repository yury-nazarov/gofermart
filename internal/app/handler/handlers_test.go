package handler

import (
	"fmt"
	"github.com/yury-nazarov/gofermart/internal/app/repository/accrual_client"
	"github.com/yury-nazarov/gofermart/internal/app/repository/pg"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yury-nazarov/gofermart/internal/app/repository/cache"
	"github.com/yury-nazarov/gofermart/internal/app/service/auth"
	"github.com/yury-nazarov/gofermart/internal/app/service/processing"
	"github.com/yury-nazarov/gofermart/internal/app/service/withdraw"
	"github.com/yury-nazarov/gofermart/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//// TODO: Временная штука пока я не научусь делать моки
////uniqUsername уникальное имя пользователя для тестов
//func uniqUsername(userName string) string{
//	bytes, err := bcrypt.GenerateFromPassword([]byte(userName), 14)
//	if err != nil {
//		log.Printf("username error: %s", err)
//	}
//	return string(bytes)
//}

// NewTestServer - конфигурируем тестовый сервер,
func NewTestServer() *httptest.Server {
	// Устанавливаем логгер
	logger := logger.NewLogger("test gofermart")

	// Иницииреуем необходимые переменные для работы сервиса из аргументов или env
	serverAddress := "127.0.0.1:8081"
	accrualAddress := "127.0.0.1"
	pgConfig := "host=localhost port=5432 user=gofermart password=gofermart dbname=gofermart sslmode=disable connect_timeout=5"

	// Инициируем БД и создаем соединение
	db := pg.NewDB(pg.DBConfig{PGConnStr: pgConfig}, logger)

	// Инициируем loginCache для проверки сессии пользователя
	// TODO: 1. Переименовать в NewLoginSession().
	//		 2. Передовать по ссылке иначе оно будет копироватся
	loginSession := cache.NewLoginCache()

	// Регистрация и авторизация пользователя
	user := auth.NewAuth(db, loginSession, logger)

	// Запускаем по тикеру горутины которые будут периодически опрашивать accrualServer и обновлять значение в БД
	accrual := httpClient.NewAccrual(accrualAddress, db, logger)

	// Бизнес логика работы с заказами
	order := processing.NewOrder(db, logger)

	// Бизнес логика работы с балансом пользователя
	balance := withdraw.NewBalance(db, logger)

	// Инициируем объект для доступа к хендлерам
	c := New(user, loginSession, order, balance, accrual, logger)

	// инициируем роутер
	router := NewRouter(c, user, logger)

	// Настраиваем адрес/порт который будут слушать тестовый сервер
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(router)
	// Закрываем созданый httptest.NewUnstartedServer Listener и назначаем подготовленный нами ранее
	// В тесткейсе нужно будет запустить и остановить сервер: ts.Start(), ts.Close()
	ts.Listener.Close()
	ts.Listener = listener
	return ts
}

// Функция HTTP клиент для тестовых запросов
func testRequest(t *testing.T, method, path string, body string, headers map[string]string) (*http.Response, string) {
	// Подготавливаем HTTP Request для тестового сервера
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	require.NoError(t, err)

	// Устанавливаем нужные хедеры для HTTP Request
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	// Убираем редирект в HTTP клиенте, для коректного тестирования HTTP хендлеров c Header Location
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)

	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestController_RegisterUser(t *testing.T) {
	// Параметры для настройки тестового HTTP Request
	type request struct {
		httpMethod string
		url        string
		headers    map[string]string
		body       string
	}
	// Ожидаемый ответ сервера
	type want struct {
		statusCode int
		headers    map[string]string
		body       string
	}
	// Список тесткейсов
	tests := []struct {
		name    string
		request request
		want    want
	}{
		// Testcases
		{
			name: "test_1: POST: Create new user [HTTP 200]",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/register",
				body:       `{ "login": "TestUser_1","password": "123123"}`,
				headers:    map[string]string{"Content-Type": "application/json"},
			},
			want: want{
				statusCode: http.StatusOK,
				//body:       ``,
				//headers:    map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "test_2: POST: User exist [HTTP 400]",
			request: request{
				// TODO: При не корректном JSON зарегал пользователя с пустыми параметрами!
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/register",
				body:       `{ "login": "TestUser_1","password": "123123"}`,
				headers:    map[string]string{"Content-Type": "application/json"},
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		//{
		//	name: "test_3: POST: Incorrect request [HTTP 409]",
		//	request: request{
		//		httpMethod: http.MethodPost,
		//		url:        "http://127.0.0.1:8081/api/user/register",
		//		body:       `{ "logi": "newUser_1","password": "104430"}`, // Incorrect JSON struct
		//		headers:    map[string]string{"Content-Type": "application/json"},
		//	},
		//	want: want{
		//		statusCode: http.StatusBadRequest,
		//	},
		//},
		{
			name: "test_3_1: POST: Incorrect request [HTTP 409]",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/register",
				body:       `{ "login": "TestUser_1","password": "123123"`, // Incorrect JSON struct
				headers:    map[string]string{"Content-Type": "application/json"},
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "test_4: POST: User exist [HTTP 409]",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/register",
				body:       `{ "login": "TestUser_1","password": "123123"}`,
				headers:    map[string]string{"Content-Type": "application/json"},
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		// TODO: How test HTTP 500 Error?
	}
	ts := NewTestServer()
	ts.Start()
	for _, tt := range tests {
		testName := fmt.Sprintf("%s.", tt.name)
		t.Run(testName, func(t *testing.T) {
			// Выполняем тестовый HTTP Request
			resp, body := testRequest(t, tt.request.httpMethod, tt.request.url, tt.request.body, tt.request.headers)
			defer resp.Body.Close() // go vet test

			assert.Equal(t, tt.want.headers["Content-Type"], resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.body, body)

		})
	}
	ts.Close()
}

// TestController_LogInUser
func TestController_LogInUser(t *testing.T) {
	// Параметры для настройки тестового HTTP Request
	type request struct {
		httpMethod string
		url        string
		headers    map[string]string
		body       string
	}
	// Ожидаемый ответ сервера
	type want struct {
		statusCode int
		headers    map[string]string
		body       string
	}
	// Список тесткейсов
	tests := []struct {
		name    string
		request request
		want    want
	}{
		// Testcases
		{
			name: "test_1: POST: LogIn [HTTP 200]",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/login",
				body:       `{ "login": "TestUser_1","password": "123123"}`,
				headers:    map[string]string{"Content-Type": "application/json"},
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "test_2: POST: Incorrect request format [HTTP 400]",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/login",
				body:       `{ "login": "TestUser_1","password": "123123"`,
				headers:    map[string]string{"Content-Type": "application/json"},
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "test_3_1: POST: Incorrect login or password [HTTP 401]",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/login",
				body:       `{ "login": "TestUser_1","password": "123"}`,
				headers:    map[string]string{"Content-Type": "application/json"},
			},
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		// TODO: How test HTTP 500 Error?
	}
	ts := NewTestServer()
	ts.Start()
	for _, tt := range tests {
		testName := fmt.Sprintf("%s.", tt.name)
		t.Run(testName, func(t *testing.T) {
			// Выполняем тестовый HTTP Request
			resp, body := testRequest(t, tt.request.httpMethod, tt.request.url, tt.request.body, tt.request.headers)
			defer resp.Body.Close() // go vet test

			assert.Equal(t, tt.want.headers["Content-Type"], resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.body, body)

		})
	}
	ts.Close()
}

// TestController_NewOrders
func TestController_NewOrders(t *testing.T) {
	// Параметры для настройки тестового HTTP Request
	type request struct {
		httpMethod string
		url        string
		headers    map[string]string
		body       string
	}
	// Ожидаемый ответ сервера
	type want struct {
		statusCode int
		headers    map[string]string
		body       string
	}
	// Список тесткейсов
	tests := []struct {
		name    string
		request request
		want    want
	}{
		// Testcases
		{
			name: "test_1: POST: Add Order ok202",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/order",
				body:       `12345678903`,
				headers:    map[string]string{"Content-Type": "text/plain"},
			},
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name: "test_2: POST: Add Order ok200. The order exist for this user",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/order",
				body:       `12345678903`,
				headers:    map[string]string{"Content-Type": "text/plain"},
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "test_3: POST: Add Order err400. Invalid request format",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/order",
				body:       `[-123-`,
				headers:    map[string]string{"Content-Type": "text/plain"},
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "test_4: POST: Add Order err422. Invalid request format",
			request: request{
				httpMethod: http.MethodPost,
				url:        "http://127.0.0.1:8081/api/user/order",
				body:       `asdasdasda`,
				headers:    map[string]string{"Content-Type": "text/plain"},
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		//// TODO: 409
		//{
		//	name: "test_4_1: POST: Add Order err409. The order has been upload other user",
		//	request: request{
		//		httpMethod: http.MethodPost,
		//		url:        "http://127.0.0.1:8081/api/user/order",
		//		body:       `12345678901`,
		//		headers:    map[string]string{
		//					"Content-Type": "text/plain",
		//					"Authorization": "xxx",
		//					},
		//	},
		//	want: want{
		//		statusCode: http.StatusBadRequest,
		//	},
		//},
		//{
		//	name: "test_4_2: POST: Add Order err409. The order has been upload other user",
		//	request: request{
		//		httpMethod: http.MethodPost,
		//		url:        "http://127.0.0.1:8081/api/user/order",
		//		body:       `123-`,
		//		headers:    map[string]string{"Content-Type": "text/plain"},
		//	},
		//	want: want{
		//		statusCode: http.StatusBadRequest,
		//	},
		//},

	}
	ts := NewTestServer()
	ts.Start()
	for _, tt := range tests {
		testName := fmt.Sprintf("%s.", tt.name)
		t.Run(testName, func(t *testing.T) {
			// Выполняем тестовый HTTP Request
			resp, body := testRequest(t, tt.request.httpMethod, tt.request.url, tt.request.body, tt.request.headers)
			defer resp.Body.Close() // go vet test

			assert.Equal(t, tt.want.headers["Content-Type"], resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.body, body)

		})
	}
	ts.Close()
}
