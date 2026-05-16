package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	dto "github.com/juliannGabrielDev/intelfy-api/internal/dto"
	"github.com/juliannGabrielDev/intelfy-api/internal/service"
	"github.com/juliannGabrielDev/intelfy-api/pkg/render"
)

type SongHandler struct {
	service *service.SongService
}

func NewSongHandler(s *service.SongService) *SongHandler {
	return &SongHandler{service: s}
}

func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	song, err := h.service.CreateSong(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, song)
}

func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	filters := dto.SongFilters{
		AlbumID:  query.Get("album-id"),
		ArtistID: query.Get("artist-id"),
		GenreID:  query.Get("genre-id"),
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	songs, err := h.service.GetSongs(r.Context(), filters)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, songs)
}

func (h *SongHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")

	if songID == "" {
		render.Error(w, http.StatusBadRequest, "song id is required")
		return
	}

	song, err := h.service.GetSongByID(r.Context(), songID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, song)
}

func (h *SongHandler) UpdateSongByID(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateSongRequest
	songID := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if songID == "" {
		render.Error(w, http.StatusBadRequest, "song id is required")
		return
	}

	if err := h.service.UpdateSongByID(r.Context(), songID, req); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Song updated successfully"})
}

func (h *SongHandler) DeleteSongByID(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")

	if songID == "" {
		render.Error(w, http.StatusBadRequest, "song id is required")
		return
	}

	if err := h.service.DeleteSongByID(r.Context(), songID); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Song deleted successfully"})
}
