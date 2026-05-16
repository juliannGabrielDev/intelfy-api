-- name: CreateAlbum :one
INSERT INTO albums (
    id, name, description, cover_url, artist_id, release_date
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetAlbums :many
SELECT alb.* FROM albums alb
WHERE (alb.artist_id = $1 OR $1 = '')
ORDER BY alb.release_date DESC
LIMIT $2 OFFSET $3;

-- name: CountAlbums :one
SELECT COUNT(*) FROM albums alb
WHERE (alb.artist_id = $1 OR $1 = '');

-- name: GetAlbumByID :one
SELECT * FROM albums
WHERE id = $1 LIMIT 1;

-- name: UpdateAlbumByID :exec
UPDATE albums
SET
    name = COALESCE(NULLIF(sqlc.arg(name), ''), name),
    description = COALESCE(NULLIF(sqlc.arg(description), ''), description),
    cover_url = COALESCE(NULLIF(sqlc.arg(cover_url), ''), cover_url),
    artist_id = COALESCE(NULLIF(sqlc.arg(artist_id), ''), artist_id),
    release_date = COALESCE(sqlc.narg(release_date), release_date)
WHERE id = sqlc.arg(id);

-- name: DeleteAlbumByID :exec
DELETE FROM albums
WHERE id = $1;
