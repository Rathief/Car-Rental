package handler

import (
	"car-rental/entity"
	"car-rental/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (rh RentalHandler) GetUserRents(c echo.Context) error {
	// get user id
	claims, err := utils.DecodeToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, err, "Error reading token")
		return err
	}
	userID := claims["userID"]

	// get records
	var records []entity.Record
	result := rh.DB.Where("user_id = ?", userID).Find(&records)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving data")
		return result.Error
	}
	c.JSON(http.StatusOK, records)
	return nil
}
func (rh RentalHandler) RentAProduct(c echo.Context) error {
	// get user id from token
	claims, err := utils.DecodeToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, err, "Error reading token")
		return err
	}
	userID := claims["userID"]
	// get user
	var user entity.User
	result := rh.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error retrieving user data")
		return result.Error
	}

	// read input
	var input entity.Rent
	if err := c.Bind(&input); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error reading input")
		return err
	}
	// get product from input
	var product entity.Product
	result = rh.DB.Where("id = ?", input.ProductID).First(&product)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving product data")
		return result.Error
	}

	// compare (rentPrice * rentLength) with userDeposit
	totalPrice := product.RentalPrice * float64(input.RentLength)
	// deny if total price > deposit
	if totalPrice > user.Deposit {
		err = fmt.Errorf("total price %.2f is larger than user deposit %.2f", totalPrice, user.Deposit)
		utils.HandleError(c, http.StatusBadRequest, err, "Not enough deposit")
		return err
	}
	tx := rh.DB.Begin()
	// create record
	result = tx.Create(&entity.Record{
		UserID:    uint(userID.(float64)),
		ProductID: product.ID,
		EndDate:   time.Now().AddDate(0, 0, int(input.RentLength)),
	})
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error inserting record")
		tx.Rollback()
		return result.Error
	}
	var record entity.Record
	result.First(&record)

	// update user by subtracting total price from deposit
	result = tx.Model(&user).Update("Deposit", user.Deposit-totalPrice)
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error updating user")
		tx.Rollback()
		return result.Error
	}
	result.First(&user)

	result = tx.Commit()
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "commit error?")
		return result.Error
	}

	if err := c.JSON(http.StatusOK, map[string]any{
		"rental_record": record,
		"user_balance":  user.Deposit,
	}); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error writing json response")
		return err
	}
	return nil
}
