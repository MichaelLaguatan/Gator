-- name: GetFeedFollowsForUser :many
SELECT feeds.name
FROM feeds
INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1;