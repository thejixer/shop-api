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

type ChargeBalanceDto struct {
	Amount float64 `json:"amount" validate:"required"`
}

type PointDto struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type CreateAddressDto struct {
	Title         string  `json:"title" validate:"required"`
	Lon           float64 `json:"lon" validate:"required"`
	Lat           float64 `json:"lat" validate:"required"`
	Address       string  `json:"address" validate:"required"`
	RecieverName  string  `json:"recieverName" validate:"required"`
	RecieverPhone string  `json:"recieverPhone" validate:"required"`
}

type CheckOutDto struct {
	AddressId int `json:"addressId" validate:"required"`
}
