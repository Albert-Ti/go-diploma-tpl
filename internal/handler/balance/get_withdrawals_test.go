package balance_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/balance"
	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository/mocks"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetWithdrawals(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	svc := service.NewService(mockRepo)
	handler := balance.GetWithdrawals(svc)

	now := time.Now()

	var results []models.WithdrawalsResp

	results = append(results, models.WithdrawalsResp{
		Order:       "9278923470",
		Sum:         751,
		ProcessedAt: now,
	})

	JSONResults, _ := json.Marshal(results)

	tests := []struct {
		name           string
		contentType    string
		responseBody   []byte
		expectedStatus int
		setupMock      func(mock *mocks.MockRepository)
	}{
		{
			name:           "Success",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			responseBody:   JSONResults,
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().GetWithdrawals(gomock.Any(), 1).
					Return(results, nil)
			},
		},
		{
			name:           "No_Content",
			expectedStatus: http.StatusNoContent,
			setupMock: func(mock *mocks.MockRepository) {
				mock.EXPECT().GetWithdrawals(gomock.Any(), 1).
					Return([]models.WithdrawalsResp{}, nil)
			},
		},
		{
			name:           "Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			setupMock:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			var req *http.Request
			if tt.expectedStatus == http.StatusUnauthorized {
				req = httptest.NewRequest(http.MethodPost, "/api/user/balance", nil)
			} else {
				req = httptest.NewRequestWithContext(
					context.WithValue(context.Background(), auth.UserIDKey, "1"),
					http.MethodGet, "/api/user/balance", nil)
			}

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if tt.expectedStatus == http.StatusOK {
				body, _ := io.ReadAll(rr.Body)
				assert.JSONEq(t, string(tt.responseBody), string(body))
			}

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}

}
