package models

type ResponseDTO struct {
	Msg        string `json:"msg"`
	StatusCode int    `json:"statusCode"`
}

type SignUpDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateAdminDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RequestVerificationEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestChangePasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type TokenDTO struct {
	Token string `json:"token"`
}

type ChangePasswordDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type CreateProductDto struct {
	Title       string  `json:"title" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	Description string  `json:"description" validate:"required"`
}

type AddtoCartDto struct {
	ProductId int `json:"productId"`
	Quantity  int `json:"quantity"`
}
