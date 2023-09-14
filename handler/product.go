package handler

import (
	"car-rental/entity"
	"car-rental/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

func (ph ProductHandler) ReadByID(c echo.Context) error {
	id := c.Param("id")

	var product entity.Product
	result := ph.DB.Preload("Records").Where("id = ?", id).First(&product)
	if result.Error != nil {
		utils.HandleError(c, http.StatusInternalServerError, result.Error, "Error retrieving data")
		return result.Error
	}
	c.JSON(http.StatusOK, product)
	return nil
}
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
	c.JSON(http.StatusAccepted, product)
	return nil
}
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

// must also delete related records
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
