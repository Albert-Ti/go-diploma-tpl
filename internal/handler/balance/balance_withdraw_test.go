package balance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/balance"
	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository/mocks"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBalanceWithdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	svc := service.NewService(mockRepo)
	handler := balance.BalanceWithdraw(svc)

	tests := []struct {
		name           string
		contentType    string
		requestBody    models.WithdrawRequest
		responseBody   models.BalanceResp
		expectedStatus int
		ctx            context.Context
		setupMock      func(mock *mocks.MockRepository)
	}{
		{
			name:           "Success",
			contentType:    "application/json",
			requestBody:    models.WithdrawRequest{Order: "9278923470", Sum: 751},
			expectedStatus: http.StatusOK,
			ctx:            context.WithValue(context.Background(), auth.UserIDKey, "1"),
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().GetBalance(gomock.Any(), 1).
					Return(models.BalanceResp{
						Current:  1000,
						Withdraw: 0,
					}, nil)

				mock.EXPECT().BalanceWithdrawalsTx(gomock.Any(), "9278923470", gomock.Any(), 1).
					Return(nil)
			},
		},
		{
			name:           "Incorrect_order_number_format",
			contentType:    "application/json",
			requestBody:    models.WithdrawRequest{Order: "1234", Sum: 0},
			expectedStatus: http.StatusUnprocessableEntity,
			ctx:            context.WithValue(context.Background(), auth.UserIDKey, "1"),
			setupMock:      nil,
		},
		{
			name:           "Insufficient_funds",
			contentType:    "application/json",
			requestBody:    models.WithdrawRequest{Order: "9278923470", Sum: 751},
			expectedStatus: http.StatusPaymentRequired,
			ctx:            context.WithValue(context.Background(), auth.UserIDKey, "1"),
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().GetBalance(gomock.Any(), 1).
					Return(models.BalanceResp{
						Current:  700,
						Withdraw: 0,
					}, nil)
			},
		},
		{
			name:           "Unauthorized",
			contentType:    "application/json",
			requestBody:    models.WithdrawRequest{Order: "9278923470", Sum: 751},
			expectedStatus: http.StatusUnauthorized,
			ctx:            context.Background(),
			setupMock:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			body, _ := json.Marshal(tt.requestBody)

			req := httptest.NewRequestWithContext(tt.ctx, http.MethodPost,
				"/api/user/balance/withdraw", bytes.NewReader(body))

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
