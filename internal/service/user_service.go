package service

import (
	"context"
	"errors"
	"math"

	"github.com/juliannGabrielDev/intelfy-api/internal/dto"
	"github.com/juliannGabrielDev/intelfy-api/internal/repository"
	"github.com/juliannGabrielDev/intelfy-api/pkg/nanoid"
	"github.com/juliannGabrielDev/intelfy-api/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.Queries
}

func NewUserService(repo *repository.Queries) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	newID, err := nanoid.GenerateID()
	if err != nil {
		return dto.AuthResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		ID:           newID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         dto.RoleUser,
	})
	if err != nil {
		return dto.AuthResponse{}, err
	}
	accessToken, err := token.GenerateToken(user.ID, user.Role.String)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token: accessToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role.String,
			CreatedAt: user.CreatedAt.Time,
		},
	}, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid credentials")
	}

	accessToken, err := token.GenerateToken(user.ID, user.Role.String)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token: accessToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role.String,
			CreatedAt: user.CreatedAt.Time,
		},
	}, nil
}

func (s *UserService) GetUsers(ctx context.Context, pagination dto.GenericPagination) (*dto.PaginatedResponse[dto.UserResponse], error) {
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	total, err := s.repo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dto.PaginatedResponse[dto.UserResponse]{
			Data: []dto.UserResponse{},
			Meta: dto.PaginationMeta{
				TotalRecords: 0,
				CurrentPage:  1,
				TotalPages:   0,
				Limit:        int(pagination.Limit),
			},
		}, nil
	}

	users, err := s.repo.GetUsers(ctx, repository.GetUsersParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})
	if err != nil {
		return nil, err
	}

	resData := make([]dto.UserResponse, len(users))
	for i, u := range users {
		resData[i] = dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role.String,
			CreatedAt: u.CreatedAt.Time,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	currentPage := int(pagination.Offset/pagination.Limit) + 1

	return &dto.PaginatedResponse[dto.UserResponse]{
		Data: resData,
		Meta: dto.PaginationMeta{
			TotalRecords: total,
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			Limit:        int(pagination.Limit),
		},
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (dto.UserResponse, error) {
	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.String,
		CreatedAt: u.CreatedAt.Time,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	existing, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	username := existing.Username
	if req.Username != "" {
		username = req.Username
	}

	email := existing.Email
	if req.Email != "" {
		email = req.Email
	}

	passwordHash := existing.PasswordHash
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return dto.UserResponse{}, err
		}
		passwordHash = string(hashed)
	}

	updated, err := s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         existing.Role,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        updated.ID,
		Username:  updated.Username,
		Email:     updated.Email,
		Role:      updated.Role.String,
		CreatedAt: updated.CreatedAt.Time,
	}, nil
}
