package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository/mocks"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	svc := service.NewService(mockRepo)
	handler := auth.Login(svc)

	const fixedSalt = "testsalt123"
	validHashedPassword := utils.HashPassword(fixedSalt, "pass")

	tests := []struct {
		name           string
		contentType    string
		body           models.AuthRequest
		expectedStatus int
		setupMock      func(mock *mocks.MockRepository)
	}{
		{
			name:           "Case_1_Success",
			contentType:    "application/json",
			body:           models.AuthRequest{Login: "test", Password: "pass"},
			expectedStatus: http.StatusOK,
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().
					GetUser(gomock.Any(), "test").
					Return(1, validHashedPassword, nil)
			},
		},
		{
			name:           "Case_2_Invalid_pass",
			contentType:    "application/json",
			body:           models.AuthRequest{Login: "test", Password: "invalidPass"},
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().
					GetUser(gomock.Any(), "test").
					Return(1, validHashedPassword, nil)
			},
		},
		{
			name:           "Case_3_Bad_Request",
			contentType:    "text",
			body:           models.AuthRequest{Login: "test", Password: "pass"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
