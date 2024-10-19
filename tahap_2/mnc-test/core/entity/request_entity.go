package entity

type RequestCreateUser struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name"  validate:"required"`
	PhoneNumber string `json:"phone_number"  validate:"required"`
	PIN         string `json:"pin"  validate:"required"`
	Address     string `json:"address"  validate:"required"`
}

type RequestUpdateUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Address   string `json:"address"  validate:"required"`
}
type RequestLoginUser struct {
	PhoneNumber string `json:"phone_number"  validate:"required"`
	PIN         string `json:"pin"  validate:"required"`
}

type RequestCreatePayment struct {
	Amount  float64 `json:"amount"  validate:"required"`
	Remarks string  `json:"remarks"  validate:"required"`
}

type RequestCreateTopUp struct {
	Amount  float64 `json:"amount"  validate:"required"`
	Remarks string  `json:"remarks"  validate:"required"`
}
type RequestCreateTransfer struct {
	TargetUser string  `json:"target_user"`
	Amount     float64 `json:"amount"  validate:"required"`
	Remarks    string  `json:"remarks"  validate:"required"`
}
