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
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
		render.Error(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	artistID := r.FormValue("artist_id")
	releaseDateStr := r.FormValue("release_date")

	var releaseDate dto.Date
	if releaseDateStr != "" {
		t, _ := time.Parse("2006-01-02", releaseDateStr)
		releaseDate = dto.Date{Time: t}
	}

	req := dto.CreateAlbumRequest{
		Name:        name,
		Description: description,
		ArtistID:    artistID,
		ReleaseDate: releaseDate,
	}

	// Handle Cover File
	var coverPath string
	file, header, err := r.FormFile("cover")
	if err == nil {
		defer file.Close()
		coverPath, err = h.saveFile(file, "albums/covers", header.Filename)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, "Failed to save cover file")
			return
		}
	}

	album, err := h.service.CreateAlbum(r.Context(), req, coverPath)
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
	albumID := chi.URLParam(r, "id")
	if albumID == "" {
		render.Error(w, http.StatusBadRequest, "album id is required")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		render.Error(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	var req dto.UpdateAlbumRequest
	if name := r.FormValue("name"); name != "" {
		req.Name = &name
	}
	if desc := r.FormValue("description"); desc != "" {
		req.Description = &desc
	}
	if artistID := r.FormValue("artist_id"); artistID != "" {
		req.ArtistID = &artistID
	}
	if relDate := r.FormValue("release_date"); relDate != "" {
		t, _ := time.Parse("2006-01-02", relDate)
		req.ReleaseDate = &dto.Date{Time: t}
	}

	var coverPathPtr *string
	file, header, err := r.FormFile("cover")
	if err == nil {
		defer file.Close()
		coverPath, err := h.saveFile(file, "albums/covers", header.Filename)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, "Failed to save cover file")
			return
		}
		coverPathPtr = &coverPath
	}

	if err := h.service.UpdateAlbumByID(r.Context(), albumID, req, coverPathPtr); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Album updated successfully"})
}

func (h *AlbumHandler) saveFile(file io.Reader, subDir, fileName string) (string, error) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	targetDir := filepath.Join(uploadDir, subDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

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

	return filepath.ToSlash(relPath), nil
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
