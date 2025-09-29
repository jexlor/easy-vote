package api

import (
	"fmt"
	"net/http"

	"github.com/jexlor/votingapp/db/store"
	"github.com/labstack/echo/v4"
)

func (s *Config) HandlerReactComment(c echo.Context) error {
	userID, ok := c.Get("userID").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req struct {
		CommentID int32 `json:"comment_id" form:"comment_id"`
		Reaction  int16 `json:"reaction" form:"reaction"` // 1=like, -1=dislike
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Reaction != 1 && req.Reaction != -1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid reaction"})
	}

	err := s.DB.CreateCommentReaction(c.Request().Context(), store.CreateCommentReactionParams{
		CommentID: req.CommentID,
		UserID:    userID,
		Reaction:  req.Reaction,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save reaction"})
	}

	counts, err := s.DB.GetCommentReactionsCount(c.Request().Context(), req.CommentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch counts"})
	}

	html := fmt.Sprintf(`<div id="comment-%d-reactions">%d 👍 %d 👎</div>`,
		req.CommentID, counts.Likes, counts.Dislikes)

	return c.HTML(http.StatusOK, html)
}
