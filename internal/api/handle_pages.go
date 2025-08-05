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
	templ.Handler(components.LoginForm()).ServeHTTP(c.Response(), c.Request())
	return nil
}

func HandleRegisternPage(c echo.Context) error {
	templ.Handler(components.RegisterForm()).ServeHTTP(c.Response(), c.Request())
	return nil
}
