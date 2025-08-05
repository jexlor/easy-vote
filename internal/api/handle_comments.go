package api

import (
	"net/http"

	"github.com/jexlor/votingapp/db/store"
	"github.com/labstack/echo/v4"
)

func (s *Config) HandlerGetAllComments(c echo.Context) error {
	comments, err := s.DB.GetAllComments(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve comments"})
	}

	return c.JSON(http.StatusOK, comments)
}

func (s *Config) HandlerCreateComment(c echo.Context) error {
	userID, ok := c.Get("userID").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req struct {
		Comment string `json:"comment"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	i, err := s.DB.CreateComment(c.Request().Context(), store.CreateCommentParams{
		UserID:  userID,
		Comment: req.Comment,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, i)
}
