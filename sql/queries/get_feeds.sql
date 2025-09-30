-- name: GetFeeds :many
SELECT users.name, feeds.name, feeds.url
FROM users
INNER JOIN feeds
ON users.id = feeds.user_id;