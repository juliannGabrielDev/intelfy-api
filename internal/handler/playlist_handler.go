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

type PlaylistHandler struct {
	service *service.PlaylistService
}

func NewPlaylistHandler(s *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service: s}
}

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePlaylistRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	playlist, err := h.service.CreatePlaylist(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, playlist)
}

func (h *PlaylistHandler) GetPlaylistsByUser(w http.ResponseWriter, r *http.Request) {
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

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		render.Error(w, http.StatusBadRequest, "user id is required")
		return
	}

	playlist, err := h.service.GetPlaylistsByUser(r.Context(), dto.GetPlaylistsRequest{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, playlist)
}

func (h *PlaylistHandler) GetPlaylistByID(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "id")
	if playlistID == "" {
		render.Error(w, http.StatusBadRequest, "playlist id is required")
		return
	}

	playlist, err := h.service.GetPlaylistByID(r.Context(), playlistID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, playlist)
}

func (h *PlaylistHandler) AddSongToPlaylist(w http.ResponseWriter, r *http.Request) {
	var req dto.AddSongToPlaylistRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.service.AddSongToPlaylist(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, result)
}

func (h *PlaylistHandler) ClearPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "id")
	if playlistID == "" {
		render.Error(w, http.StatusBadRequest, "playlist id is required")
		return
	}

	if err := h.service.ClearPlaylist(r.Context(), playlistID); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, "Playlist cleared successfully")
}

func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "id")
	if playlistID == "" {
		render.Error(w, http.StatusBadRequest, "playlist id is required")
		return
	}

	playlist, err := h.service.DeletePlaylist(r.Context(), playlistID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, playlist)
}
