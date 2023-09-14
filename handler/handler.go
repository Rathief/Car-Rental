package handler

import "gorm.io/gorm"

type UserHandler struct {
	DB *gorm.DB
}
type ProductHandler struct {
	DB *gorm.DB
}
type RentalHandler struct {
	DB *gorm.DB
}
