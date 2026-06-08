package models

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetOrdersResp struct {
	Number     string `json:"number"`
	Status     int    `json:"status"`
	Accrual    string `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}
