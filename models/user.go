package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserRepository interface {
	Create(data SignUpDTO) (*User, error)
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	VerifyEmail(email string) error
	UpdatePassword(email, password string) error
}

type User struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Email           string          `json:"email"`
	IsEmailVerified bool            `json:"isEmailVerified"`
	Password        string          `json:"password"`
	Balance         decimal.Decimal `json:"balance"`
	CreatedAt       time.Time       `json:"createdAt"`
}

type UserDto struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}
