package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthUserID(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantID  int
		wantErr bool
	}{
		{
			name:    "Success",
			ctx:     context.WithValue(context.Background(), auth.UserIDKey, "123"),
			wantID:  123,
			wantErr: false,
		},
		{
			name:    "Empty userID",
			ctx:     context.WithValue(context.Background(), auth.UserIDKey, ""),
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "No userID in context",
			ctx:     context.Background(),
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "Invalid userID (not string)",
			ctx:     context.WithValue(context.Background(), auth.UserIDKey, 123),
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := auth.GetAuthUserID(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}

func TestGuard(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetAuthUserID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(userID)))
	})

	protectedHandler := auth.Guard(testHandler)

	tests := []struct {
		name           string
		setupCookie    func() *http.Cookie
		expectedStatus int
	}{
		{
			name: "Valid token",
			setupCookie: func() *http.Cookie {
				token, _ := auth.CreateToken("123")
				return auth.CreateCookie("token", token)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing cookie",
			setupCookie: func() *http.Cookie {
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token",
			setupCookie: func() *http.Cookie {
				return &http.Cookie{Name: "token", Value: "invalid-token"}
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.setupCookie != nil {
				if cookie := tt.setupCookie(); cookie != nil {
					req.AddCookie(cookie)
				}
			}

			rr := httptest.NewRecorder()
			protectedHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
