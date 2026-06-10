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

func Register(svc *service.Service) http.HandlerFunc {
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

		userID, err := svc.Register(r.Context(), req.Login, req.Password)
		if err != nil {
			if errors.Is(err, service.ErrLoginAlreadyExists) {
				http.Error(w, err.Error(), http.StatusConflict)
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
