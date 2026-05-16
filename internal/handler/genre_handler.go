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

type GenreHandler struct {
	service *service.GenreService
}

func NewGenreHandler(s *service.GenreService) *GenreHandler {
	return &GenreHandler{service: s}
}

func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	genre, err := h.service.CreateGenre(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, genre)
}

func (h *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
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

	genres, err := h.service.GetGenres(r.Context(), pagination)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, genres)
}

func (h *GenreHandler) GetGenreByID(w http.ResponseWriter, r *http.Request) {
	genreID := chi.URLParam(r, "id")

	if genreID == "" {
		render.Error(w, http.StatusBadRequest, "genre id is required")
		return
	}

	genre, err := h.service.GetGenreByID(r.Context(), genreID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, genre)
}

func (h *GenreHandler) UpdateGenreByID(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateGenreRequest
	genreID := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if genreID == "" {
		render.Error(w, http.StatusBadRequest, "genre id is required")
		return
	}

	genre, err := h.service.UpdateGenreByID(r.Context(), genreID, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, genre)
}

func (h *GenreHandler) DeleteGenreByID(w http.ResponseWriter, r *http.Request) {
	genreID := chi.URLParam(r, "id")

	if genreID == "" {
		render.Error(w, http.StatusBadRequest, "genre id is required")
		return
	}

	genre, err := h.service.DeleteGenreByID(r.Context(), genreID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, genre)
}
