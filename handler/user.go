package handler

import (
	"car-rental/entity"
	"car-rental/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (uh UserHandler) RegisterUser(c echo.Context) error {
	// read input
	var user entity.User
	if err := c.Bind(&user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error Reading JSON Input")
		return err
	}

	// hash pass
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error hashing pass")
		return err
	}
	user.Password = string(hashedPass)

	// insert data
	result := uh.DB.Select("name", "email", "password").Create(&user)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error inserting")
		return result.Error
	}
	result.First(&user)

	// generate and send token
	err = utils.GenerateToken(c, user.ID, user.Role)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error generating token")
		return err
	}
	return nil
}

func (uh UserHandler) LoginUser(c echo.Context) error {
	var user entity.User
	var storedUser entity.User
	if err := c.Bind(&user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error Reading JSON Input")
		return err
	}
	result := uh.DB.Where("email = ?", user.Email).First(&storedUser)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving data")
		return result.Error
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Invalid Password")
		return err
	}

	// generate and send token
	err = utils.GenerateToken(c, storedUser.ID, storedUser.Role)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error generating token")
		return err
	}
	return nil
}

func (uh UserHandler) ReadAll(c echo.Context) error {
	var users []entity.User
	result := uh.DB.Preload("Records").Find(&users)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving data")
		return result.Error
	}
	c.JSON(http.StatusOK, users)
	return nil
}

func (uh UserHandler) ReadByID(c echo.Context) error {
	id := c.Param("id")

	var user entity.User
	result := uh.DB.Preload("Records").Where("id = ?", id).First(&user)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving data")
		return result.Error
	}
	c.JSON(http.StatusOK, user)
	return nil
}

func (uh UserHandler) TopUpDeposit(c echo.Context) error {
	// get input
	var topUp entity.TopUp
	if err := c.Bind(&topUp); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error reading input")
		return err
	}

	// get user from auth token
	var user entity.User
	claims, err := utils.DecodeToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error reading token")
		return err
	}
	id := claims["userID"]
	result := uh.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving data")
		return result.Error
	}

	// add input deposit to user deposit
	user.Deposit += topUp.Deposit
	result = uh.DB.Save(&user)
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error updating data")
		return result.Error
	}
	c.JSON(http.StatusOK, map[string]any{
		"Current Deposit": user.Deposit,
	})
	return nil
}
