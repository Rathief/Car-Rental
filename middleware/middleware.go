package middleware

import (
	"car-rental/utils"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.DecodeToken(c)
		if err != nil {
			utils.HandleError(c, http.StatusUnauthorized, err, "Authorization error")
			return nil
		}
		if claims.Valid() != nil {
			utils.HandleError(c, http.StatusUnauthorized, claims.Valid(), "Unauthorized user")
			return nil
		}
		return next(c)
	}
}
func AuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.DecodeToken(c)
		if err != nil {
			utils.HandleError(c, http.StatusUnauthorized, err, "Authorization error")
			return nil
		}
		if claims.Valid() != nil {
			utils.HandleError(c, http.StatusUnauthorized, claims.Valid(), "Unauthorized user")
			return nil
		}
		if claims["userRole"] != "admin" {
			utils.HandleError(c, http.StatusUnauthorized, fmt.Errorf("user is not an admin"), "Unauthorized user")
			return nil
		}
		return next(c)
	}
}

func Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := getSession(c)
		if session.Values["authenticated"] == nil {
			session.Values["authenticated"] = false
		}

		if err := session.Save(c.Request(), c.Response()); err != nil {
			return err
		}

		return next(c)
	}
}

func getSession(c echo.Context) *sessions.Session {
	session, _ := utils.CookieStore.Get(c.Request(), "session-test")
	return session
}
