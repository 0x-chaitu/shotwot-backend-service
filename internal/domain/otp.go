package domain

type Otp struct {
	Phone   string `json:"phone"`
	Otp     string `json:"otp"`
	OrderId string `json:"orderId"`
}
