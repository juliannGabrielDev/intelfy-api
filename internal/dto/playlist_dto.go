package dto

import "time"

type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
	UserID      string `json:"user_id" binding:"required"`
}

type GetPlaylistsRequest struct {
	UserID string `url:"user_id" binding:"required"`
	Limit  int32  `url:"limit"`
	Offset int32  `url:"offset"`
}

type AddSongToPlaylistRequest struct {
	PlaylistID string `json:"playlist_id" binding:"required"`
	SongID     string `json:"song_id" binding:"required"`
}

type PlaylistResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetPlaylistResponse struct {
	Playlist PlaylistResponse `json:"playlist"`
	Songs    []SongResponse   `json:"songs"`
}

type AddSongToPlaylistResponse struct {
	PlaylistID string       `json:"playlist_id"`
	TotalSongs int          `json:"total_songs"`
	AddedSong  SongResponse `json:"added_song"`
}
