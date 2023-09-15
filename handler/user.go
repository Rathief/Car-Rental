package handler

import (
	"car-rental/entity"
	"car-rental/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser godoc
//
//	@Summary		Register User
//	@Description	Register a user by json, notify the registered account, and returns a jwt token. Email will be validated first.
//	@Tags			User
//	@Accept			json
//	@Produce		json,html
//	@Param			account	body		entity.User	true	"Register user"
//	@Success		201		{object}	entity.User
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Router			/users/register [post]
func (uh UserHandler) RegisterUser(c echo.Context) error {
	// read input
	var user entity.User
	if err := c.Bind(&user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error Reading JSON Input")
		return err
	}

	// validate email
	emailURL := strings.Replace(user.Email, "@", "%40", -1)
	url := fmt.Sprintf("https://email-validator28.p.rapidapi.com/email-validator/validate?email=%s", emailURL)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", os.Getenv("EMAIL_VALID_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("EMAIL_VALID_HOST"))

	res, _ := http.DefaultClient.Do(req)
	var emailVal entity.EmailValidate
	err := json.NewDecoder(res.Body).Decode(&emailVal)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error using email validator")
		return err
	}
	if !emailVal.IsValid || !emailVal.IsDeliverable {
		err = fmt.Errorf("email not valid or deliverable")
		utils.HandleError(c, http.StatusBadRequest, err, fmt.Sprintf("isValid = %v, isDeliverable = %v", emailVal.IsValid, emailVal.IsDeliverable))
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

	// generate token
	token, err := utils.GenerateToken(c, user.ID, user.Role)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error generating token, try logging in.")
		return err
	}
	// show token
	err = c.JSON(http.StatusCreated, map[string]any{
		"token": token,
	})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error displaying token, try logging in.")
		return err
	}
	// send email notification
	err = utils.SendEmail(user.Email, fmt.Sprintf("Welcome to Car Rental, %s!", user.Name), "<h1>Welcome!</h1><br><p>You have successfully registered to Car Rental.</p>")
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error sending email")
		return err
	}
	return nil
}

// LoginUser godoc
//
//	@Summary		Login User
//	@Description	Login by json and returns jwt token
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			account	body		entity.User	true	"Login user"
//	@Success		201		{object}	string
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Router			/users/login [post]
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
	token, err := utils.GenerateToken(c, storedUser.ID, storedUser.Role)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error generating token")
		return err
	}
	err = c.JSON(http.StatusCreated, map[string]any{
		"token": token,
	})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error displaying token, try logging in again.")
		return err
	}

	return nil
}

// ReadAll godoc
//
//	@Summary		Show all users
//	@Description	Show all users and their rents in JSON form
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.User
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		401	{object}	utils.ErrorResponse
//	@Router			/users/ [get]
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

// ReadByID godoc
//
//	@Summary		Show user
//	@Description	Show user by id from url
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	entity.User
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		401	{object}	utils.ErrorResponse
//	@Router			/users/{id} [get]
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

// TopUpDeposit godoc
//
//	@Summary		Top up user deposit
//	@Description	Top up the user's deposit by the specified amount, and send an email notification
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			amount	body		entity.TopUp	true	"Top up amount"
//	@Success		200		{object}	entity.TopUp	"User deposit after top up"
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		401		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Router			/users/topup [post]
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

	// send email notification
	err = utils.SendEmail(user.Email, "Top Up Successful!", fmt.Sprintf("<h1>Top Up Successful!</h1><br><p>Your Car Rental Deposit is now %v.</p>", user.Deposit))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Error sending email")
		return err
	}
	return nil
}
