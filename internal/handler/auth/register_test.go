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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	svc := service.NewService(mockRepo)
	handler := auth.Register(svc)

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
				mock.EXPECT().RegisterTx(gomock.Any(), "test", gomock.Any()).
					Return(1, nil)
			},
		},
		{
			name:           "Case_2_Bad_Request",
			contentType:    "text",
			body:           models.AuthRequest{Login: "test", Password: "pass"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Case_3_Conflict",
			contentType:    "application/json",
			body:           models.AuthRequest{Login: "test", Password: "pass"},
			expectedStatus: http.StatusConflict,
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().
					RegisterTx(gomock.Any(), "test", gomock.Any()).
					Return(0, service.ErrLoginAlreadyExists)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
