package models

import (
	"time"
)

type AuthRequest struct {
	Login    string `json:"login" example:"Albert"`
	Password string `json:"password" example:"12345"`
}

type WithdrawRequest struct {
	Order string  `json:"order" example:"9278923470"`
	Sum   float64 `json:"sum" example:"751"`
}

type OrdersResp struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type BalanceResp struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}

type AccrualResp struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type WithdrawalsResp struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type TaskOrder struct {
	OrderID string
	UserID  int
}
