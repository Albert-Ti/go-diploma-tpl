package balance

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
)

// BalanceWithdraw godoc
//
//	@Summary		Запрос на списание средств
//	@Description	Запрос на списание средств
//	@Tags			balance
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.WithdrawRequest	true	"Формат запроса"
//
//	@Security		CookieAuth
//	@Success		200		"Успешная обработка запроса"
//	@Failure		401		"Пользователь не аутентифицирован"
//	@Failure		402		"На счету недостаточно средств"
//	@Failure		422		"Неверный номер заказа"
//	@Failure		500		"Внутренняя ошибка сервера"
//
//	@Router			/api/user/balance/withdraw [post]
func BalanceWithdraw(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m models.WithdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !utils.AlgoLuna(m.Order) {
			http.Error(w, "Incorrect order number format", http.StatusUnprocessableEntity)
			return
		}

		userID, err := auth.GetAuthUserID(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := svc.BalanceWithdrawals(r.Context(), m.Order, m.Sum, userID); err != nil {
			if errors.Is(err, service.ErrInsufficientFunds) {
				http.Error(w, err.Error(), http.StatusPaymentRequired)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
