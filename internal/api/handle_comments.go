package api

import (
	"html"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/jexlor/votingapp/db/store"
	"github.com/jexlor/votingapp/web/components"
	"github.com/jexlor/votingapp/web/models"
	"github.com/labstack/echo/v4"
)

func (s *Config) HandlerGetAllComments(c echo.Context) error {
	userID, _ := c.Get("userID").(int32)

	dbComments, err := s.DB.GetAllCommentsWithReactions(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve comments"})
	}

	comments := make([]models.CommentWithReactions, len(dbComments))
	hasCommented := false
	for i, c := range dbComments {
		comments[i] = models.CommentWithReactions{
			ID:        c.ID,
			UserID:    c.UserID,
			Comment:   c.Comment,
			CreatedAt: c.CreatedAt.Time,
			Likes:     int32(c.Likes.(int64)),
			Dislikes:  int32(c.Dislikes.(int64)),
		}
		if c.UserID == userID {
			hasCommented = true
		}

	}
	data := models.CommentsPageData{
		Comments:      comments,
		CurrentUserID: userID,
		HasCommented:  hasCommented,
	}
	templ.Handler(components.CommentsPage(data)).ServeHTTP(c.Response(), c.Request())
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

	comment = html.EscapeString(comment)

	if len(comment) > 1000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "comment too long"})
	}

	i, err := s.DB.CreateComment(c.Request().Context(), store.CreateCommentParams{
		UserID:  userID,
		Comment: comment,
	})
	if err != nil {
		return c.HTML(http.StatusOK, `<div class="error">Failed to create comment</div>`)
	}

	c.Redirect(http.StatusSeeOther, "/v1/comments")
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
