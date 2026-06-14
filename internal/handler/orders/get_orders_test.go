package orders_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/orders"
	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository/mocks"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	svc := service.NewService(mockRepo)
	handler := orders.GetOrders(svc)

	now := time.Now()

	var results []models.OrdersResp

	results = append(results, models.OrdersResp{
		Number:     "9278923470",
		Status:     models.StatusNew,
		UploadedAt: now,
	})

	JSONResults, _ := json.Marshal(results)

	tests := []struct {
		name           string
		contentType    string
		expectedStatus int
		expectedBody   []byte
		setupMock      func(mock *mocks.MockRepository)
	}{
		{
			name:           "Case_1_Success",
			expectedStatus: http.StatusOK,
			expectedBody:   JSONResults,
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().
					GetOrders(gomock.Any(), 1).
					Return(results, nil)
			},
		},
		{
			name:           "Case_2_No_Content",
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().
					GetOrders(gomock.Any(), 1).
					Return(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			req := httptest.NewRequestWithContext(
				context.WithValue(context.Background(), auth.UserIDKey, "1"),
				http.MethodGet, "/api/user/orders", nil)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if tt.expectedStatus == http.StatusOK {
				body, _ := io.ReadAll(rr.Body)
				assert.JSONEq(t, string(tt.expectedBody), string(body))
			}
			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
