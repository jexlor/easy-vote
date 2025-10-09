package api

import (
	"database/sql"
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
	sort := c.QueryParam("sort")

	var comments []models.CommentWithReactions
	hasCommented := false

	switch sort {
	case "newest":
		dbCommentsReactionsTime, err := s.DB.GetAllCommentsWithReactionsTime(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve comments"})
		}

		comments = make([]models.CommentWithReactions, len(dbCommentsReactionsTime))
		for i, r := range dbCommentsReactionsTime {
			comments[i] = models.CommentWithReactions{
				ID:        r.ID,
				UserID:    r.UserID,
				Comment:   r.Comment,
				CreatedAt: r.CreatedAt.Time,
				Likes:     int32(r.Likes.(int64)),
				Dislikes:  int32(r.Dislikes.(int64)),
			}
			if r.UserID == userID {
				hasCommented = true
			}
		}

	default: // sort=liked or none
		dbCommentsReactions, err := s.DB.GetAllCommentsWithReactions(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve comments"})
		}

		comments = make([]models.CommentWithReactions, len(dbCommentsReactions))
		for i, r := range dbCommentsReactions {
			comments[i] = models.CommentWithReactions{
				ID:        r.ID,
				UserID:    r.UserID,
				Comment:   r.Comment,
				CreatedAt: r.CreatedAt.Time,
				Likes:     int32(r.Likes.(int64)),
				Dislikes:  int32(r.Dislikes.(int64)),
			}
			if r.UserID == userID {
				hasCommented = true
			}
		}
	}

	data := models.CommentsPageData{
		Comments:      comments,
		CurrentUserID: userID,
		HasCommented:  hasCommented,
	}

	csrfToken := c.Get("csrf").(string)
	csrfTokenMap := map[string]interface{}{"CSRFToken": csrfToken}

	templ.Handler(components.CommentsPage(data, csrfTokenMap)).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (s *Config) HandlerGetCommentsBySearch(c echo.Context) error {
	userID, _ := c.Get("userID").(int32)
	query := c.QueryParam("q")
	if query == "" {
		return c.Redirect(http.StatusSeeOther, "/v1/comments")
	}

	hasCommented := false

	dbCommentsReactions, err := s.DB.GetCommentsBySearch(c.Request().Context(), sql.NullString{
		String: query,
		Valid:  true,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve comments"})
	}

	comments := make([]models.CommentWithReactions, len(dbCommentsReactions))
	for i, r := range dbCommentsReactions {
		comments[i] = models.CommentWithReactions{
			ID:        r.ID,
			UserID:    r.UserID,
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt.Time,
			Likes:     int32(r.Likes.(int64)),
			Dislikes:  int32(r.Dislikes.(int64)),
		}
		if r.UserID == userID {
			hasCommented = true
		}
	}

	data := models.CommentsPageData{
		Comments:      comments,
		CurrentUserID: userID,
		HasCommented:  hasCommented,
	}

	csrfToken := c.Get("csrf").(string)
	csrfTokenMap := map[string]interface{}{"CSRFToken": csrfToken}

	templ.Handler(components.CommentsPage(data, csrfTokenMap)).ServeHTTP(c.Response(), c.Request())
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
