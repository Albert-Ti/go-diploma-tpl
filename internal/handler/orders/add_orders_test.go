package orders_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/orders"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository/mocks"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	svc := service.NewService(mockRepo)

	wp := orders.NewWorkerPool(svc, 0, 10)
	defer wp.Stop()

	handler := orders.AddOrders(svc, wp)

	tests := []struct {
		name           string
		contentType    string
		requestBody    string
		expectedStatus int
		ctx            context.Context
		setupMock      func(mock *mocks.MockRepository)
	}{
		{
			name:           "Success",
			contentType:    "application/json",
			requestBody:    "9278923470",
			expectedStatus: http.StatusAccepted,
			ctx:            context.WithValue(context.Background(), auth.UserIDKey, "1"),
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().CreateOrder(gomock.Any(), "9278923470", 1).
					Return(0, nil)
			},
		},
		{
			name:           "Order_already_exists_for_this_user",
			contentType:    "application/json",
			requestBody:    "9278923470",
			expectedStatus: http.StatusOK,
			ctx:            context.WithValue(context.Background(), auth.UserIDKey, "1"),
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().CreateOrder(gomock.Any(), "9278923470", 1).
					Return(0, service.ErrOrderAlreadyExistsForUser)
			},
		},
		{
			name:           "Order_already_exists_for_other_user",
			contentType:    "application/json",
			requestBody:    "9278923470",
			expectedStatus: http.StatusConflict,
			ctx:            context.WithValue(context.Background(), auth.UserIDKey, "2"),
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().CreateOrder(gomock.Any(), "9278923470", 2).
					Return(1, service.ErrOrderAlreadyExistsForOther)
			},
		},
		{
			name:           "Incorrect_order_number_format",
			contentType:    "application/json",
			requestBody:    "12345",
			expectedStatus: http.StatusUnprocessableEntity,
			ctx:            context.WithValue(context.Background(), auth.UserIDKey, "1"),
			setupMock:      nil,
		},
		{
			name:           "Unauthorized",
			contentType:    "application/json",
			requestBody:    "9278923470",
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

			req := httptest.NewRequestWithContext(tt.ctx, http.MethodPost,
				"/api/user/balance/withdraw", strings.NewReader(tt.requestBody))

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}

}
