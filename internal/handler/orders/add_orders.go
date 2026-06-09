package orders

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
)

func AddOrders(svc *service.Service, wp *WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, _ := io.ReadAll(r.Body)
		orderID := string(body)

		if !utils.AlgoLuna(orderID) {
			http.Error(w, "Incorrect order number format", http.StatusUnprocessableEntity)
			return
		}

		userID, err := auth.GetAuthUserID(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		errOrder := svc.AddOrder(r.Context(), orderID, id)

		switch {
		case errors.Is(errOrder, service.ErrOrderAlreadyExistsForUser):
			http.Error(w, errOrder.Error(), http.StatusOK)
		case errors.Is(errOrder, service.ErrOrderAlreadyExistsForOther):
			http.Error(w, errOrder.Error(), http.StatusConflict)
		case errOrder == nil:
			w.WriteHeader(http.StatusAccepted)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		hash, _ := utils.RandomHash(5)

		wp.AddTask(hash)

		go func() {
			slog.Info("results", "response", <-wp.results)

		}()
	}
}
