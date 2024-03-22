package domain

type Otp struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Otp     string `json:"otp"`
	OrderId string `json:"orderId"`
}
