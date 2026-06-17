package balance

import (
	"encoding/json"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
)

// GetBalance godoc
//
//	@Summary		 Получение текущего баланса пользователя
//	@Description	Получение текущего баланса пользователя
//	@Tags				balance
//	@Produce		json
//
//	@Security		CookieAuth
//	@Success		200		"успешная обработка запроса"
//	@Failure		401		"Пользователь не аутентифицирован"
//	@Failure		500		"Внутренняя ошибка сервера"
//
//	@Router			/api/user/balance [get]
func GetBalance(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetAuthUserID(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		balance, err := svc.GetBalance(r.Context(), userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(balance); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
