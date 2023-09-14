package utils

import "github.com/labstack/echo/v4"

type ErrorResponse struct {
	Error   string
	Details any
}

func HandleError(c echo.Context, status int, err error, details any) {
	c.JSON(status, ErrorResponse{
		Error:   err.Error(),
		Details: details,
	})
}
