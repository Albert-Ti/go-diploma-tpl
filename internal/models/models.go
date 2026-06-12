package models

import (
	"time"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetOrdersResp struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type GetBalanceResp struct {
	Current  float64 `json:"current"`
	Withdraw int     `json:"withdrawn"`
}

type AccrualResp struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type TaskOrder struct {
	OrderID string
	UserID  int
}
