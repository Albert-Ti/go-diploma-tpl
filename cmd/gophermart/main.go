package main

import (
	"log/slog"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {

	repo, e := repository.NewRepository()
	if e != nil {
		panic(e)
	}

	svc := service.NewService(repo)

	r := chi.NewRouter()

	r.Post("/api/user/register", handler.Register(svc))
	r.Post("/api/user/login", handler.Login(svc))
	r.Post("/api/user/orders", handler.AddOrders(svc))
	r.Post("/api/user/balance/withdraw", handler.BalanceWithdraw(svc))

	r.Get("/api/user/orders", handler.GetOrders(svc))
	r.Get("/api/user/balance", handler.GetBalance(svc))
	r.Get("/api/user/withdrawals", handler.StatusBalance(svc))

	slog.Info("Running server", "host", config.Envs.RunAddr)

	err := http.ListenAndServe(config.Envs.RunAddr, r)

	if err != nil {
		panic(err)
	}
}
