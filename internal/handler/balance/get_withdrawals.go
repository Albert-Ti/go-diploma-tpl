package balance

import (
	"encoding/json"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
)

// GetWithdrawals godoc
//
//	@Summary		 Мои заказы
//	@Description	Получение списка загруженных номеров заказов
//	@Tags				balance
//	@Produce		json
//
//	@Security		CookieAuth
//	@Success		200		"Успешная обработка запроса"
//	@Failure		204		"Hет ни одного списания"
//	@Failure		401		"Пользователь не аутентифицирован"
//	@Failure		500		"Внутренняя ошибка сервера"
//
//	@Router			/api/user/withdrawals [get]
func GetWithdrawals(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetAuthUserID(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		withdrawals, err := svc.GetWithdrawals(r.Context(), userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
