-- name: GetPosts :many
SELECT posts.*, feeds.name AS feed_name
FROM posts
INNER JOIN feeds ON feeds.id = posts.feed_id
INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;