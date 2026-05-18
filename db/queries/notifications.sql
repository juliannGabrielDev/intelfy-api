-- name: CreateNotification :exec
INSERT INTO notifications (id, user_id, title, message)
VALUES ($1, $2, $3, $4);

-- name: GetNotificationsByUserID :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: MarkAsRead :exec
UPDATE notifications
SET is_read = TRUE
WHERE id = $1;
