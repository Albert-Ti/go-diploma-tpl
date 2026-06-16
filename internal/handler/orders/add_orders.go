package orders

import (
	"errors"
	"io"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
)

// AddOrders godoc
//
//	@Summary		 Создать заказ
//	@Description	Создание заказа а так же проверка в системе лояльности
//	@Tags				orders
//	@Accept			text/plain
//	@Param			request	body string true "Номер заказа" example(9278923470)
//
//	@Security		CookieAuth
//	@Success		200		"Номер заказа уже был загружен этим пользователем"
//	@Success		202		"Новый номер заказа принят в обработку"
//	@Failure		400		"Неверный формат запроса"
//	@Failure		401		"Пользователь не аутентифицирован"
//	@Failure		409		"Номер заказа уже был загружен другим пользователем"
//	@Failure		422		"Неверный формат номера заказа"
//	@Failure		500		"Внутренняя ошибка сервера"
//
//	@Router			/api/user/orders [post]
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

		errOrder := svc.AddOrder(r.Context(), orderID, userID)

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

		wp.AddTask(orderID, userID)
	}
}
