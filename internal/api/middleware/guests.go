package middleware

// import (
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// )

// func GuestOnlyMiddleware() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			userID, ok := c.Get("userID").(int32)
// 			if ok && userID != 0 {
// 				return c.Redirect(http.StatusSeeOther, "/comments")
// 			}
// 			return next(c)
// 		}
// 	}
// }
