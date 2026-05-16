-- name: CreateArtist :one
INSERT INTO artists (
	id, name, bio, cover_url
) VALUES (
	$1, $2, $3, $4
)
RETURNING id, name, bio, cover_url;

-- name: GetArtists :many
SELECT * FROM artists
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: CountArtists :one
SELECT COUNT(*) FROM artists;

-- name: GetArtistByID :one
SELECT * FROM artists
WHERE id = $1 LIMIT 1;

-- name: UpdateArtistByID :one
UPDATE artists
SET 
	name = COALESCE(NULLIF(sqlc.arg(name), ''), name),
	bio = COALESCE(NULLIF(sqlc.arg(bio), ''), bio),
	cover_url = COALESCE(NULLIF(sqlc.arg(cover_url), ''), cover_url)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteArtistByID :one
DELETE FROM artists
WHERE id = $1
RETURNING *;