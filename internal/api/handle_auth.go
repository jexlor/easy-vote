package api

import (
	"net/http"

	"github.com/jexlor/votingapp/db/store"
	"github.com/jexlor/votingapp/internal/auth"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Config) HandleLogin(c echo.Context) error {
	var req struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	user, err := s.DB.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}
	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/v1/comments")
	return c.NoContent(http.StatusOK)

}

func (s *Config) HandleRegister(c echo.Context) error {
	var req struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	_, err = s.DB.CreateUser(c.Request().Context(), store.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return c.JSON(http.StatusConflict, err)
	}
	c.Response().Header().Set("HX-Redirect", "/v1/login")
	return c.NoContent(http.StatusOK)
}
