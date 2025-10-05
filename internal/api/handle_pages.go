package api

import (
	"github.com/a-h/templ"
	"github.com/jexlor/votingapp/web/components"
	"github.com/labstack/echo/v4"
)

func HandleLoginPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	data := map[string]interface{}{
		"CSRFToken": csrfToken,
	}
	templ.Handler(components.LoginPage(data)).ServeHTTP(c.Response(), c.Request())
	return nil
}
func HandleRegisterPage(c echo.Context) error {
	templ.Handler(components.RegisterPage()).ServeHTTP(c.Response(), c.Request())
	return nil
}
