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

type AlbumHandler struct {
	service *service.AlbumService
}

func NewAlbumHandler(s *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{service: s}
}

func (h *AlbumHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	album, err := h.service.CreateAlbum(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, album)
}

func (h *AlbumHandler) GetAlbums(w http.ResponseWriter, r *http.Request) {
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

	filters := dto.AlbumFilters{
		ArtistID: query.Get("artist-id"),
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	albums, err := h.service.GetAlbums(r.Context(), filters)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, albums)
}

func (h *AlbumHandler) GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "id")

	if albumID == "" {
		render.Error(w, http.StatusBadRequest, "album id is required")
		return
	}

	album, err := h.service.GetAlbumByID(r.Context(), albumID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, album)
}

func (h *AlbumHandler) UpdateAlbumByID(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateAlbumRequest
	albumID := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if albumID == "" {
		render.Error(w, http.StatusBadRequest, "album id is required")
		return
	}

	if err := h.service.UpdateAlbumByID(r.Context(), albumID, req); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Album updated successfully"})
}

func (h *AlbumHandler) DeleteAlbumByID(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "id")

	if albumID == "" {
		render.Error(w, http.StatusBadRequest, "album id is required")
		return
	}

	if err := h.service.DeleteAlbumByID(r.Context(), albumID); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Album deleted successfully"})
}
