package main

import (
	"log/slog"
	"net/http"

	_ "github.com/Albert-Ti/go-diploma-tpl/docs"
	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/balance"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/orders"
	"github.com/Albert-Ti/go-diploma-tpl/internal/middleware"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Gophermart API
// @version 1.0
// @description Сервис лояльности Gophermart
// @host localhost:8080
// @BasePath /
func main() {
	config.ParseFlag()

	repo, e := repository.NewRepository()
	if e != nil {
		panic(e)
	}

	svc := service.NewService(repo)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.GzipCompress)

	r.Post("/api/user/register", auth.Register(svc))
	r.Post("/api/user/login", auth.Login(svc))

	wp := orders.NewWorkerPool(svc, 3, 10)
	defer wp.Stop()

	r.Post("/api/user/orders", auth.Guard(orders.AddOrders(svc, wp)))
	r.Get("/api/user/orders", auth.Guard(orders.GetOrders(svc)))
	r.Post("/api/user/balance/withdraw", auth.Guard(balance.BalanceWithdraw(svc)))
	r.Get("/api/user/balance", auth.Guard(balance.GetBalance(svc)))
	r.Get("/api/user/withdrawals", auth.Guard(balance.GetWithdrawals(svc)))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	slog.Info("Running server", "host", config.Envs.RunAddr)

	err := http.ListenAndServe(config.Envs.RunAddr, r)

	if err != nil {
		panic(err)
	}
}
