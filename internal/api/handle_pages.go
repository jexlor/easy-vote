package api

import (
	"github.com/a-h/templ"
	"github.com/jexlor/votingapp/web/components"
	"github.com/labstack/echo/v4"
)

func HandleHomePage(c echo.Context) error {
	templ.Handler(components.HomePage()).ServeHTTP(c.Response(), c.Request())
	return nil
}

func HandleLoginPage(c echo.Context) error {
	templ.Handler(components.LoginPage()).ServeHTTP(c.Response(), c.Request())
	return nil
}
func HandleRegisterPage(c echo.Context) error {
	templ.Handler(components.RegisterPage()).ServeHTTP(c.Response(), c.Request())
	return nil
}
