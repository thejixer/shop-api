package models

type AddressRepository interface {
	Create(a *Address) (*Address, error)
	FindById(id int) (*Address, error)
	FindByUserId(userId int) ([]*Address, error)
	Delete(id int) error
}

type Address struct {
	Id            int     `json:"id"`
	UserId        int     `json:"userId"`
	Title         string  `json:"title"`
	Lon           float64 `json:"lon"`
	Lat           float64 `json:"lat"`
	Address       string  `json:"address"`
	RecieverName  string  `json:"recieverName"`
	RecieverPhone string  `json:"recieverPhone"`
	Deleted       bool    `json:"deleted"`
}

type AddressDto struct {
	Id            int      `json:"id"`
	User          UserDto  `json:"user"`
	Title         string   `json:"title"`
	Coordinates   PointDto `json:"coordinates"`
	Address       string   `json:"address"`
	RecieverName  string   `json:"recieverName"`
	RecieverPhone string   `json:"recieverPhone"`
}
