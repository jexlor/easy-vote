package api

import (
	"net/http"
	"regexp"
	"unicode"

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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isStrongPassword(pw string) bool {
	if len(pw) < 8 || len(pw) > 64 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func (s *Config) HandleRegister(c echo.Context) error {
	var req struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if !emailRegex.MatchString(req.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid email format"})
	}

	if !isStrongPassword(req.Password) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password must be 8 chars, include uppercase, lowercase, and a number",
		})
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
		return c.JSON(http.StatusConflict, map[string]string{"error": "email already registered"})
	}

	c.Response().Header().Set("HX-Redirect", "/v1/login")
	return c.NoContent(http.StatusOK)
}
