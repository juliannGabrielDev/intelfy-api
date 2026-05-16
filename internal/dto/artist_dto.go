package dto

type CreateArtistRequest struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Bio      string `json:"bio,omitempty"`
	CoverURL string `json:"cover_url,omitempty"`
}

type UpdateArtistRequest struct {
	Name     *string `json:"name"`
	Bio      *string `json:"bio,omitempty"`
	CoverURL *string `json:"cover_url,omitempty"`
}

type ArtistResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Bio      string `json:"bio,omitempty"`
	CoverURL string `json:"cover_url,omitempty"`
}
