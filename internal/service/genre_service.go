package service

import (
	"context"
	"errors"
	"math"

	"github.com/juliannGabrielDev/intelfy-api/internal/dto"
	"github.com/juliannGabrielDev/intelfy-api/internal/repository"
	"github.com/juliannGabrielDev/intelfy-api/pkg/nanoid"
)

type GenreService struct {
	repo *repository.Queries
}

func NewGenreService(repo *repository.Queries) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) CreateGenre(ctx context.Context, req dto.CreateGenreRequest) (dto.GenreResponse, error) {
	if req.Name == "" {
		return dto.GenreResponse{}, errors.New("genre name cannot be empty")
	}

	newId, err := nanoid.GenerateID()
	if err != nil {
		return dto.GenreResponse{}, err
	}

	genre, err := s.repo.CreateGenre(ctx, repository.CreateGenreParams{
		ID:   newId,
		Name: req.Name,
	})

	if err != nil {
		return dto.GenreResponse{}, err
	}

	return dto.GenreResponse{ID: genre.ID, Name: genre.Name}, nil
}

func (s *GenreService) GetGenres(ctx context.Context, pagination dto.GenericPagination) (*dto.PaginatedResponse[dto.GenreResponse], error) {
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	total, err := s.repo.CountGenres(ctx)

	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dto.PaginatedResponse[dto.GenreResponse]{
			Data: []dto.GenreResponse{},
			Meta: dto.PaginationMeta{
				TotalRecords: 0,
				CurrentPage:  1,
				TotalPages:   1,
				Limit:        int(pagination.Limit),
			},
		}, nil
	}

	genres, err := s.repo.GetGenres(ctx, repository.GetGenresParams{Limit: pagination.Limit, Offset: pagination.Offset})

	if err != nil {
		return nil, err
	}

	resData := make([]dto.GenreResponse, len(genres))
	for i, genre := range genres {
		resData[i] = dto.GenreResponse{ID: genre.ID, Name: genre.Name}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	currentPage := int(pagination.Offset/pagination.Limit) + 1

	return &dto.PaginatedResponse[dto.GenreResponse]{
		Data: resData,
		Meta: dto.PaginationMeta{
			TotalRecords: total,
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			Limit:        int(pagination.Limit),
		},
	}, nil
}

func (s *GenreService) GetGenreByID(ctx context.Context, id string) (dto.GenreResponse, error) {
	genre, err := s.repo.GetGenreByID(ctx, id)

	if err != nil {
		return dto.GenreResponse{}, err
	}

	return dto.GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}, nil
}

func (s *GenreService) UpdateGenreByID(ctx context.Context, id string, req dto.UpdateGenreRequest) (dto.GenreResponse, error) {
	if req.Name != nil && *req.Name == "" {
		return dto.GenreResponse{}, errors.New("genre name cannot be empty")
	}

	genre, err := s.repo.UpdateGenreByID(ctx, repository.UpdateGenreByIDParams{
		Name: req.Name,
		ID:   id,
	})
	if err != nil {
		return dto.GenreResponse{}, err
	}

	return dto.GenreResponse{ID: genre.ID, Name: genre.Name}, nil
}

func (s *GenreService) DeleteGenreByID(ctx context.Context, id string) (dto.GenreResponse, error) {
	genre, err := s.repo.DeleteGenreByID(ctx, id)
	if err != nil {
		return dto.GenreResponse{}, err
	}

	return dto.GenreResponse{ID: genre.ID, Name: genre.Name}, nil
}
