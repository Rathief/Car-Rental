package utils

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var CookieStore = sessions.NewCookieStore([]byte("test"))

func getSession(c echo.Context) *sessions.Session {
	session, _ := CookieStore.Get(c.Request(), "session-test")
	return session
}
