package handler

import (
	"car-rental/entity"
	"car-rental/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ReadAll godoc
//
//	@Summary		Show all products
//	@Description	Show all products and related rents
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.Product
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		401	{object}	utils.ErrorResponse
//	@Failure		500	{object}	utils.ErrorResponse
//	@Router			/products/ [get]
func (ph ProductHandler) ReadAll(c echo.Context) error {
	var products []entity.Product
	result := ph.DB.Preload("Records").Find(&products)
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error retrieving data")
		return result.Error
	}
	c.JSON(http.StatusOK, products)
	return nil
}

// ReadByID godoc
//
//	@Summary		Show product
//	@Description	Show product by id from url
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	entity.Product
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		401	{object}	utils.ErrorResponse
//	@Router			/products/{id} [get]
func (ph ProductHandler) ReadByID(c echo.Context) error {
	id := c.Param("id")

	var product entity.Product
	result := ph.DB.Preload("Records").Where("id = ?", id).First(&product)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving data")
		return result.Error
	}
	c.JSON(http.StatusOK, product)
	return nil
}

// CreateProduct godoc
//
//	@Summary		Create product
//	@Description	Insert new product data
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			product	body		entity.Product	true	"Product Data"
//	@Success		201		{object}	entity.Product
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		401		{object}	utils.ErrorResponse
//	@Router			/products/ [post]
func (ph ProductHandler) CreateProduct(c echo.Context) error {
	// get input
	var product entity.Product
	if err := c.Bind(&product); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error reading input")
		return err
	}

	// insert data
	result := ph.DB.Create(&product)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error inserting data")
		return result.Error
	}
	result.Preload("Records").First(&product)
	c.JSON(http.StatusCreated, product)
	return nil
}

// UpdateProduct godoc
//
//	@Summary		Update product
//	@Description	Update product targeted by the given ID using given product data
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Product ID"
//	@Param			product	body		entity.Product	true	"Product Data"
//	@Success		200		{object}	entity.Product
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		401		{object}	utils.ErrorResponse
//	@Router			/products/{id} [put]
func (ph ProductHandler) UpdateProductByID(c echo.Context) error {
	// get input
	var product entity.Product
	if err := c.Bind(&product); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Error reading input")
		return err
	}

	// get id from param
	id := c.Param("id")

	// update data
	result := ph.DB.Model(&entity.Product{}).Where("id = ?", id).Updates(product)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error updating data")
		return result.Error
	}
	result.Preload("Records").First(&product)
	c.JSON(http.StatusAccepted, product)
	return nil
}

// DeleteProduct godoc
//
//	@Summary		Delete product
//	@Description	Delete product targeted by the given ID and related rent data
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	string
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		401	{object}	utils.ErrorResponse
//	@Failure		500	{object}	utils.ErrorResponse
//	@Router			/products/{id} [delete]
func (ph ProductHandler) DeleteProductByID(c echo.Context) error {
	// get id from param
	productID := c.Param("id")

	// get product by ID
	var product entity.Product
	result := ph.DB.Preload("Records").Where("id = ?", productID).First(&product)
	if result.Error != nil {
		utils.HandleError(c, http.StatusBadRequest, result.Error, "Error retrieving product data")
		return result.Error
	}

	tx := ph.DB.Begin()
	for _, record := range product.Records {
		result := tx.Delete(record)
		if result.Error != nil {
			utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error deleting related records")
			tx.Rollback()
			return result.Error
		}
	}
	result = tx.Delete(product)
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error deleting product")
		tx.Rollback()
		return result.Error
	}
	result = tx.Commit()
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "commit error?")
		return result.Error
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "product successfully deleted",
	})
	return nil
}
