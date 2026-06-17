package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
)

// Login godoc
//
//	@Summary		Авторизация пользователя
//	@Description	Авторизует пользователя и устанавливает JWT cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.AuthRequest	true	"Данные пользователя"
//
//	@Success		200		"Пользователь успешно зарегистрирован"
//	@Failure		400		"Некорректный запрос"
//	@Failure		401		"Неправильный логин или пароль"
//	@Failure		500		"Внутренняя ошибка сервера"
//
//	@Router			/api/user/login [post]
func Login(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-type"), "application/json") {
			http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
			return
		}

		var req models.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID, err := svc.Login(r.Context(), req.Login, req.Password)
		if err != nil {
			if errors.Is(err, service.ErrUnauthorized) {
				http.Error(w, "Invalid login or password", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		token, err := CreateToken(strconv.Itoa(userID))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, CreateCookie("token", token))

		w.WriteHeader(http.StatusOK)
	}
}
