package main

import (
	"car-rental/config"
	"car-rental/handler"
	"car-rental/middleware"
	"log"

	_ "car-rental/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			Car Rental API
//	@version		0.1
//	@description	A Car Rental API for H8 Phase 2 Project

//	@license.name	None

//	@host		localhost:8080
//	@BasePath	/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	db := config.ConnectDB()
	uh := handler.UserHandler{DB: db}
	ph := handler.ProductHandler{DB: db}
	rh := handler.RentalHandler{DB: db}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	u := e.Group("/users")
	u.POST("/register", uh.RegisterUser)
	u.POST("/login", uh.LoginUser)
	u.POST("/topup", uh.TopUpDeposit, middleware.Auth)
	u.GET("/", uh.ReadAll, middleware.AuthAdmin)
	u.GET("/:id", uh.ReadByID, middleware.AuthAdmin)

	p := e.Group("/products")
	p.GET("/", ph.ReadAll, middleware.Auth)
	p.GET("/:id", ph.ReadByID, middleware.Auth)
	p.POST("/", ph.CreateProduct, middleware.AuthAdmin)
	p.PUT("/:id", ph.UpdateProductByID, middleware.AuthAdmin)
	p.DELETE("/:id", ph.DeleteProductByID, middleware.AuthAdmin)

	r := e.Group("/rent")
	r.GET("/", rh.GetUserRents, middleware.Auth)
	r.POST("/", rh.RentAProduct, middleware.Auth)

	e.Logger.Fatal(e.Start(":8080"))
}
