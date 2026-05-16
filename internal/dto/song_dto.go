package dto

type CreateSongRequest struct {
	Name     string       `json:"name"`
	Duration SongDuration `json:"duration"`
	AudioURL string       `json:"audio_url"`
	AlbumID  string       `json:"album_id"`
	GenreID  string       `json:"genre_id,omitempty"`
}

type UpdateSongRequest struct {
	Name     *string       `json:"name"`
	Duration *SongDuration `json:"duration"`
	AudioURL *string       `json:"audio_url"`
	AlbumID  *string       `json:"album_id"`
	GenreID  *string       `json:"genre_id,omitempty"`
}

type SongResponse struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Duration SongDuration `json:"duration"`
	AudioURL string       `json:"audio_url"`
	AlbumID  string       `json:"album_id"`
	GenreID  string       `json:"genre_id,omitempty"`
}

type SongFilters struct {
	AlbumID  string `json:"album_id"`
	ArtistID string `json:"artist_id"`
	GenreID  string `json:"genre_id"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}
