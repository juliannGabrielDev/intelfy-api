-- name: CreateGenre :one
INSERT INTO genres (
	id, name
) VALUES (
	$1, $2
)
RETURNING *;

-- name: GetGenres :many
SELECT * FROM genres
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: CountGenres :one
SELECT COUNT(*) from genres;

-- name: GetGenreByID :one
SELECT * FROM genres
WHERE id = $1 LIMIT 1;

-- name: UpdateGenreByID :one
UPDATE genres
SET name = COALESCE(NULLIF(sqlc.arg(name), ''), name)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteGenreByID :one
DELETE FROM genres
WHERE id = sqlc.arg(id)
RETURNING *;
