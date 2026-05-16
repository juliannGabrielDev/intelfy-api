-- name: CreatePlaylist :one
INSERT INTO playlists (
	id, name, description, user_id
) VALUES (
	$1, $2, $3, $4
)
RETURNING *;

-- name: GetPlaylistsByUser :many
SELECT * FROM playlists
WHERE user_id = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: CountPlaylistsByUser :one
SELECT COUNT(*) FROM playlists
WHERE user_id = $1;

-- name: GetPlaylistByID :one
SELECT * FROM playlists
WHERE id = $1 LIMIT 1;

-- name: GetSongsByPlaylistID :many
SELECT s.* FROM songs s
JOIN playlist_songs ps ON s.id = ps.song_id
WHERE ps.playlist_id = $1
ORDER BY ps.added_at DESC;

-- name: CountSongsInPlaylist :one
SELECT COUNT(*) FROM playlist_songs
WHERE playlist_id = $1;

-- name: AddSongToPlaylist :exec
INSERT INTO playlist_songs (
	playlist_id, song_id
) VALUES (
	$1, $2
) ON CONFLICT (playlist_id, song_id) DO NOTHING;

-- name: ClearPlaylist :exec
DELETE FROM playlist_songs
WHERE playlist_id = $1;

-- name: RemoveSongFromPlaylist :exec
DELETE FROM playlist_songs
WHERE playlist_id = $1 AND song_id = $2;

-- name: DeletePlaylistByID :one
DELETE FROM playlists
WHERE id = $1
RETURNING *;
