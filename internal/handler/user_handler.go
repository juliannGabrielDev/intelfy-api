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

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authRes, err := h.service.Register(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, authRes)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authRes, err := h.service.Login(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, authRes)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.service.GetUsers(r.Context(), pagination)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		render.Error(w, http.StatusBadRequest, "user id is required")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		render.Error(w, http.StatusBadRequest, "user id is required")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), userID, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, user)
}
