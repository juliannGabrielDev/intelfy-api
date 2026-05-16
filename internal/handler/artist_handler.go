package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/juliannGabrielDev/intelfy-api/internal/dto"
	"github.com/juliannGabrielDev/intelfy-api/internal/service"
	"github.com/juliannGabrielDev/intelfy-api/pkg/render"
)

type ArtistHandler struct {
	service *service.ArtistService
}

func NewArtistHandler(s *service.ArtistService) *ArtistHandler {
	return &ArtistHandler{service: s}
}

func (h *ArtistHandler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateArtistRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	artist, err := h.service.CreateArtist(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, artist)
}

func (h *ArtistHandler) GetArtists(w http.ResponseWriter, r *http.Request) {
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

	pagination := dto.GenericPagination{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	artists, err := h.service.GetArtists(r.Context(), pagination)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artists)
}

func (h *ArtistHandler) GetArtistByID(w http.ResponseWriter, r *http.Request) {
	artistID := chi.URLParam(r, "id")
	if artistID == "" {
		render.Error(w, http.StatusBadRequest, "artist id is required")
		return
	}

	artist, err := h.service.GetArtistByID(r.Context(), artistID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artist)
}

func (h *ArtistHandler) UpdateArtistByID(w http.ResponseWriter, r *http.Request) {
	artistID := chi.URLParam(r, "id")
	if artistID == "" {
		render.Error(w, http.StatusBadRequest, "artist id is required")
		return
	}

	var req dto.UpdateArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	artist, err := h.service.UpdateArtistByID(r.Context(), artistID, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artist)
}

func (h *ArtistHandler) DeleteArtistByID(w http.ResponseWriter, r *http.Request) {
	artistID := chi.URLParam(r, "id")
	if artistID == "" {
		render.Error(w, http.StatusBadRequest, "artist id is required")
		return
	}

	artist, err := h.service.DeleteArtistByID(r.Context(), artistID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artist)
}
