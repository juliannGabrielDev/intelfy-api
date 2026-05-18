package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
	if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB
		render.Error(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	name := r.FormValue("name")
	durationStr := r.FormValue("duration")
	albumID := r.FormValue("album_id")
	genreID := r.FormValue("genre_id")

	duration, _ := time.ParseDuration(durationStr)
	if duration <= 0 {
		// fallback if it's just seconds
		if s, err := strconv.ParseFloat(durationStr, 64); err == nil {
			duration = time.Duration(s * float64(time.Second))
		}
	}

	req := dto.CreateSongRequest{
		Name:     name,
		Duration: dto.SongDuration{Duration: duration},
		AlbumID:  albumID,
		GenreID:  genreID,
	}

	// Handle Audio File
	audioFile, audioHeader, err := r.FormFile("audio")
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Audio file is required")
		return
	}
	defer audioFile.Close()

	audioPath, err := h.saveFile(audioFile, "songs/audio", audioHeader.Filename)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to save audio file")
		return
	}

	// Handle Cover File (Optional)
	var coverPath string
	coverFile, coverHeader, err := r.FormFile("cover")
	if err == nil {
		defer coverFile.Close()
		coverPath, err = h.saveFile(coverFile, "songs/covers", coverHeader.Filename)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, "Failed to save cover file")
			return
		}
	}

	song, err := h.service.CreateSong(r.Context(), req, audioPath, coverPath)
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
	songID := chi.URLParam(r, "id")
	if songID == "" {
		render.Error(w, http.StatusBadRequest, "song id is required")
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		render.Error(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	var req dto.UpdateSongRequest
	if name := r.FormValue("name"); name != "" {
		req.Name = &name
	}
	if durationStr := r.FormValue("duration"); durationStr != "" {
		duration, _ := time.ParseDuration(durationStr)
		if duration <= 0 {
			if s, err := strconv.ParseFloat(durationStr, 64); err == nil {
				duration = time.Duration(s * float64(time.Second))
			}
		}
		if duration > 0 {
			req.Duration = &dto.SongDuration{Duration: duration}
		}
	}
	if albumID := r.FormValue("album_id"); albumID != "" {
		req.AlbumID = &albumID
	}
	if genreID := r.FormValue("genre_id"); genreID != "" {
		req.GenreID = &genreID
	}

	var audioPathPtr, coverPathPtr *string

	// Handle Audio File
	audioFile, audioHeader, err := r.FormFile("audio")
	if err == nil {
		defer audioFile.Close()
		audioPath, err := h.saveFile(audioFile, "songs/audio", audioHeader.Filename)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, "Failed to save audio file")
			return
		}
		audioPathPtr = &audioPath
	}

	// Handle Cover File
	coverFile, coverHeader, err := r.FormFile("cover")
	if err == nil {
		defer coverFile.Close()
		coverPath, err := h.saveFile(coverFile, "songs/covers", coverHeader.Filename)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, "Failed to save cover file")
			return
		}
		coverPathPtr = &coverPath
	}

	if err := h.service.UpdateSongByID(r.Context(), songID, req, audioPathPtr, coverPathPtr); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Song updated successfully"})
}

func (h *SongHandler) saveFile(file io.Reader, subDir, fileName string) (string, error) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	// Create directory if it doesn't exist
	targetDir := filepath.Join(uploadDir, subDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	// Create a unique filename
	ext := filepath.Ext(fileName)
	uniqueName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	relPath := filepath.Join(subDir, uniqueName)
	fullPath := filepath.Join(uploadDir, relPath)

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	// Convert backslashes to forward slashes for URL consistency
	return filepath.ToSlash(relPath), nil
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
