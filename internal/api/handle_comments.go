package api

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/jexlor/votingapp/db/store"
	"github.com/jexlor/votingapp/web/components"
	"github.com/jexlor/votingapp/web/models"
	"github.com/labstack/echo/v4"
)

func (s *Config) HandlerGetAllComments(c echo.Context) error {
	dbComments, err := s.DB.GetAllCommentsWithReactions(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve comments"})
	}

	comments := make([]models.CommentWithReactions, len(dbComments))
	for i, c := range dbComments {
		comments[i] = models.CommentWithReactions{
			ID:        c.ID,
			UserID:    c.UserID,
			Comment:   c.Comment,
			CreatedAt: c.CreatedAt.Time,
			Likes:     int32(c.Likes.(int64)),
			Dislikes:  int32(c.Dislikes.(int64)),
		}
	}

	templ.Handler(components.CommentsPage(comments)).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (s *Config) HandlerCreateComment(c echo.Context) error {
	userID, ok := c.Get("userID").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	comment := c.FormValue("comment")
	if comment == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "empty comment"})
	}

	i, err := s.DB.CreateComment(c.Request().Context(), store.CreateCommentParams{
		UserID:  userID,
		Comment: comment,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, i)
}

func (s *Config) HandlerDeleteComment(c echo.Context) error {
	userID, ok := c.Get("userID").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Read comment ID from URL param, not form
	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid comment id"})
	}

	// Delete comment in DB
	err = s.DB.DeleteComment(c.Request().Context(), store.DeleteCommentParams{
		ID:     int32(commentID),
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete comment"})
	}

	// Redirect back to comments page
	return c.Redirect(http.StatusSeeOther, "/v1/comments")
}
