-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched = $2, updated_at = $3
WHERE feeds.id = $1;