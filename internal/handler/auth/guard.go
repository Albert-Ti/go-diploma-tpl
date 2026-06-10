package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey string = "userID"

type MyCustomClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func CreateToken(userID string) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	t.Claims = &MyCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		userID,
	}

	return t.SignedString([]byte(config.Envs.JWTSecretKey))
}

func CreateCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
	}
}

func Guard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &MyCustomClaims{}

		t, err := jwt.ParseWithClaims(
			cookie.Value,
			claims,
			func(t *jwt.Token) (any, error) {
				return []byte(config.Envs.JWTSecretKey), nil
			},
		)

		if err != nil || !t.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func GetAuthUserID(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return 0, errors.New("User ID not found")
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, errors.New("User ID not found")
	}
	return id, nil
}
