package entity

type TopUp struct {
	Deposit float64 `json:"deposit"`
}

type Rent struct {
	ProductID  uint `json:"product_id"`
	RentLength uint `json:"rent_length"`
}
