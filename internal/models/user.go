package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserRepository interface {
	Create(name, email, password, role string, isEmailVerified bool, permissions []string) (*User, error)
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	VerifyEmail(email string) error
	UpdatePassword(email, password string) error
	FindUsers(text string, page, limit int) ([]*User, int, error)
	ChargeBalance(userId int, amount float64) error
	UpdatePermissions(id int, permissions []string) error
}

type User struct {
	ID              int             `json:"id"`
	Role            string          `json:"role"`
	Name            string          `json:"name"`
	Email           string          `json:"email"`
	IsEmailVerified bool            `json:"isEmailVerified"`
	Password        string          `json:"password"`
	Balance         decimal.Decimal `json:"balance"`
	Permissions     []string        `json:"permissions"`
	CreatedAt       time.Time       `json:"createdAt"`
}

type UserDto struct {
	ID        int       `json:"id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type AdminDto struct {
	ID          int       `json:"id"`
	Role        string    `json:"role"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Balance     float64   `json:"balance"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LL_UserDto struct {
	Total  int       `json:"total"`
	Result []UserDto `json:"result"`
}

type LL_AdminDto struct {
	Total  int        `json:"total"`
	Result []AdminDto `json:"result"`
}
