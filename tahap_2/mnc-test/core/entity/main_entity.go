package entity

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenClaims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}

type TransactionType string

var (
	Payment  TransactionType = "payment"
	TopUp    TransactionType = "topUp"
	Transfer TransactionType = "transfer"
)

type User struct {
	UserID      string     `json:"user_id" sql:"type:uuid;default:uuid_generate_v4()"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `json:"phone_number"`
	PIN         string     `json:"pin"`
	Address     string     `json:"address"`
	CreatedDate *time.Time `json:"created_date,omitempty" gorm:"autoCreateTime" sql:"type:timestamptz"`
	UpdatedDate *time.Time `json:"update_date,omitempty" gorm:"autoUpdateTime" sql:"type:timestamptz"`
}

type UserBalance struct {
	UserID      string     `json:"user_id"`
	Balance     float64    `json:"balance"`
	CreatedDate *time.Time `json:"created_date,omitempty" gorm:"autoCreateTime" sql:"type:timestamptz"`
	UpdatedDate *time.Time `json:"update_date,omitempty" gorm:"autoUpdateTime" sql:"type:timestamptz"`
}

type Transaction struct {
	UserID          string          `json:"user_id"`
	TransferID      string          `json:"transfer_id"`
	PaymentID       string          `json:"payment_id"`
	TopUpID         string          `json:"top_up_id"`
	Amount          float64         `json:"amount"` // Using int64 to handle larger amounts
	Remarks         string          `json:"remarks"`
	Status          string          `json:"status"`
	TransactionType TransactionType `json:"transaction_type"`
	BalanceBefore   float64         `json:"balance_before"` // Using int64 for balance values
	BalanceAfter    float64         `json:"balance_after"`
	CreatedDate     *time.Time      `json:"created_date,omitempty" gorm:"autoCreateTime" sql:"type:timestamptz"`
	UpdatedDate     *time.Time      `json:"update_date,omitempty" gorm:"autoUpdateTime" sql:"type:timestamptz"`
}
