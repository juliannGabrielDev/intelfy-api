-- name: CreateSong :one
INSERT INTO songs (
  id, name, duration_seconds, audio_url, album_id, genre_id
) VALUES (
	$1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetSongs :many
SELECT s.id, s.name, s.duration_seconds, s.audio_url, s.album_id, s.genre_id, s.created_at FROM songs s
JOIN albums a ON s.album_id = a.id
WHERE (s.album_id = $1 OR $1 = '')
  AND (a.artist_id = $2 OR $2 = '')
  AND (s.genre_id = $3 OR $3 IS NULL)
ORDER BY s.created_at DESC
LIMIT $4 OFFSET $5;

-- name: CountSongs :one
SELECT COUNT(*) FROM songs s
JOIN albums a ON s.album_id = a.id
WHERE (s.album_id = $1 OR $1 = '')
  AND (a.artist_id = $2 OR $2 = '')
  AND (s.genre_id = $3 OR $3 IS NULL);

-- name: GetSongByID :one
SELECT id, name, duration_seconds, audio_url, album_id, genre_id, created_at FROM songs
WHERE id = $1 LIMIT 1;

-- name: UpdateSongByID :exec
UPDATE songs
SET
  name = COALESCE(NULLIF(sqlc.arg(name), ''), name),
  duration_seconds = CASE WHEN sqlc.arg(duration_seconds) > 0 THEN sqlc.arg(duration_seconds) ELSE duration_seconds END,
  audio_url = COALESCE(NULLIF(sqlc.arg(audio_url), ''), audio_url),
  album_id = COALESCE(NULLIF(sqlc.arg(album_id), ''), album_id),
  genre_id = COALESCE(NULLIF(sqlc.arg(genre_id), ''), genre_id)
WHERE id = sqlc.arg(id);

-- name: DeleteSongByID :exec
DELETE FROM songs
WHERE id = $1;
