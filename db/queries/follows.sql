-- name: FollowArtist :exec
INSERT INTO follows (follower_id, artist_id)
VALUES ($1, $2);

-- name: UnfollowArtist :exec
DELETE FROM follows
WHERE follower_id = $1 AND artist_id = $2;

-- name: GetFollowersByArtistID :many
SELECT follower_id FROM follows
WHERE artist_id = $1;
