package orders

import (
	"encoding/json"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
)

// GetOrders godoc
//
//	@Summary		 Мои заказы
//	@Description	Получение списка загруженных номеров заказов
//	@Tags				orders
//	@Produce		json
//
//	@Security		CookieAuth
//	@Success		200		"Успешная обработка запроса"
//	@Failure		204		"Нет данных для ответа"
//	@Failure		401		"Пользователь не аутентифицирован"
//	@Failure		500		"Внутренняя ошибка сервера"
//
//	@Router			/api/user/orders [get]
func GetOrders(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetAuthUserID(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		orders, err := svc.GetOrders(r.Context(), userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
