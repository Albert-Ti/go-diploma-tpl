package main

import (
	"log/slog"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/balance"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/orders"
	"github.com/Albert-Ti/go-diploma-tpl/internal/middleware"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

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

	r.Post("/api/user/orders", auth.Guard(orders.AddOrders(svc)))
	r.Get("/api/user/orders", auth.Guard(orders.GetOrders(svc)))
	r.Post("/api/user/balance/withdraw", auth.Guard(balance.BalanceWithdraw(svc)))
	r.Get("/api/user/balance", auth.Guard(balance.GetBalance(svc)))
	r.Get("/api/user/withdrawals", auth.Guard(balance.StatusBalance(svc)))

	slog.Info("Running server", "host", config.Envs.RunAddr)

	err := http.ListenAndServe(config.Envs.RunAddr, r)

	if err != nil {
		panic(err)
	}
}
