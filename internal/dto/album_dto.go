package dto

import "time"

type CreateAlbumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CoverURL    string `json:"cover_url,omitempty"`
	ArtistID    string `json:"artist_id"`
	ReleaseDate Date   `json:"release_date"`
}

type UpdateAlbumRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description,omitempty"`
	CoverURL    *string `json:"cover_url,omitempty"`
	ArtistID    *string `json:"artist_id"`
	ReleaseDate *Date   `json:"release_date"`
}

type AlbumResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CoverURL    string    `json:"cover_url,omitempty"`
	ArtistID    string    `json:"artist_id"`
	ReleaseDate time.Time `json:"release_date"`
}

type AlbumFilters struct {
	ArtistID string `json:"artist_id"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}
