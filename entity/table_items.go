package entity

import "time"

type User struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique;"`
	Password string  `json:"password"`
	Deposit  float64 `json:"deposit" gorm:"default:0"`
	Role     string  `json:"role" gorm:"default:customer"` // customer,admin
	Records  []Record
}
type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	RentalPrice float64 `json:"rental_price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"` // car,motorcycle
	Records     []Record
}
type Record struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	StartDate time.Time `json:"start_date" gorm:"autoCreateTime"`
	EndDate   time.Time `json:"end_date"`
}
