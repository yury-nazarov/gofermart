# The temporary information about the project

Postgres запущен в контейнере:
```shell
cd deployments
docker-compose up
```

Run app with arguments
```shell
go run cmd/gophermart/main.go -a "127.0.0.1:8080" -r "127.0.0.1" -d "host=localhost port=5432 user=gofermart password=gofermart dbname=gofermart sslmode=disable connect_timeout=5"
```

Run app with EnvVars

```

```

## Тесты

```shell
go test
go test -cover
go test -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
open cover.html
```

## Progress

- [ ] Init
  - [x] Run app with args
  - [x] Run app with env
  - [x] Logs subsustem
  - [ ] Middleware
    - [ ] compress
    - [ ] Auth Check
- [ ] Handlers
  - [ ] POST /api/user/register
    - [ ] Tests
  - [ ] POST /api/user/login
    - [ ] Tests
  - [ ] POST /api/user/orders
    - [ ] Tests
  - [ ] GET /api/user/orders
    - [ ] Tests
  - [ ] GET /api/user/balance
    - [ ] Tests
  - [ ] POST /api/user/balance/withdraw
    - [ ] Tests
  - [ ] GET /api/user/balance/withdrawals
    - [ ] Tests