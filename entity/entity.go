package entity

type TopUp struct {
	Deposit float64 `json:"deposit"`
}

type Rent struct {
	ProductID  uint `json:"product_id"`
	RentLength uint `json:"rent_length"`
}

type EmailValidate struct {
	OriginalMail  string `json:"originalMail"`
	Message       string `json:"message"`
	IsDisposable  bool   `json:"isDisposable"`
	IsValid       bool   `json:"isValid"`
	IsDeliverable bool   `json:"isDeliverable"`
}
