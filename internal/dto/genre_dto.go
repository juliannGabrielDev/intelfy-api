package dto

type CreateGenreRequest struct {
	Name string `json:"name"`
}

type UpdateGenreRequest struct {
	Name *string `json:"name"`
}

type GenreResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
