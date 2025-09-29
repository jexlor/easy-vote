-- name: CreateCommentReaction :exec
INSERT INTO comment_reactions (comment_id, user_id, reaction)
VALUES ($1, $2, $3)
ON CONFLICT (comment_id, user_id)
DO UPDATE SET reaction = EXCLUDED.reaction;

-- name: GetAllCommentsWithReactions :many
SELECT
    c.id,
    c.user_id,
    c.comment,
    c.created_at,
    COALESCE(SUM((CASE WHEN r.reaction = 1 THEN 1 ELSE 0 END)::int4), 0) AS likes,
    COALESCE(SUM((CASE WHEN r.reaction = -1 THEN 1 ELSE 0 END)::int4), 0) AS dislikes
FROM comments c
LEFT JOIN comment_reactions r ON r.comment_id = c.id
GROUP BY c.id
ORDER BY likes DESC, c.created_at DESC;


-- name: GetCommentReactionsCount :one
SELECT
  COALESCE(SUM(CASE WHEN reaction = 1 THEN 1 ELSE 0 END), 0) AS likes,
  COALESCE(SUM(CASE WHEN reaction = -1 THEN 1 ELSE 0 END), 0) AS dislikes
FROM comment_reactions
WHERE comment_id = $1;



