package models

import "time"

type CommentWithReactions struct {
	ID        int32
	UserID    int32
	Comment   string
	CreatedAt time.Time
	Likes     int32
	Dislikes  int32
}

type GetCommentReactionsCountRow struct {
	Likes    int32
	Dislikes int32
}
