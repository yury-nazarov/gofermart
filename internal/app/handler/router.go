package handler

import (
	"github.com/yury-nazarov/gofermart/internal/app/repository/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter - создает роутер. Внутри определяем все ручки
//				TODO: user auth.UserInterface - для передачи в Middleware во время регистрации
func NewRouter(c *Controller, user auth.UserInterface) http.Handler {
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
