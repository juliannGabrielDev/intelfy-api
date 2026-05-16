package service

import (
	"context"
	"errors"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juliannGabrielDev/intelfy-api/internal/dto"
	"github.com/juliannGabrielDev/intelfy-api/internal/repository"
)

type ArtistService struct {
	repo *repository.Queries
}

func NewArtistService(repo *repository.Queries) *ArtistService {
	return &ArtistService{repo: repo}
}

func (s *ArtistService) CreateArtist(ctx context.Context, req dto.CreateArtistRequest) (dto.ArtistResponse, error) {
	if req.UserID == "" {
		return dto.ArtistResponse{}, errors.New("user_id is required")
	}

	artist, err := s.repo.CreateArtist(ctx, repository.CreateArtistParams{
		ID:   req.UserID,
		Name: req.Name,
		Bio: pgtype.Text{
			String: req.Bio,
			Valid:  req.Bio != "",
		},
		CoverUrl: pgtype.Text{
			String: req.CoverURL,
			Valid:  req.CoverURL != "",
		},
	})
	if err != nil {
		return dto.ArtistResponse{}, err
	}

	return dto.ArtistResponse{
		ID:       artist.ID,
		Name:     artist.Name,
		Bio:      artist.Bio.String,
		CoverURL: artist.CoverUrl.String,
	}, nil
}

func (s *ArtistService) GetArtists(ctx context.Context, pagination dto.GenericPagination) (*dto.PaginatedResponse[dto.ArtistResponse], error) {
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	total, err := s.repo.CountArtists(ctx)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dto.PaginatedResponse[dto.ArtistResponse]{
			Data: []dto.ArtistResponse{},
			Meta: dto.PaginationMeta{
				TotalRecords: 0,
				CurrentPage:  1,
				TotalPages:   0,
				Limit:        int(pagination.Limit),
			},
		}, nil
	}

	artists, err := s.repo.GetArtists(ctx, repository.GetArtistsParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})
	if err != nil {
		return nil, err
	}

	resData := make([]dto.ArtistResponse, len(artists))
	for i, artist := range artists {
		resData[i] = dto.ArtistResponse{
			ID:       artist.ID,
			Name:     artist.Name,
			Bio:      artist.Bio.String,
			CoverURL: artist.CoverUrl.String,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	currentPage := int(pagination.Offset/pagination.Limit) + 1

	return &dto.PaginatedResponse[dto.ArtistResponse]{
		Data: resData,
		Meta: dto.PaginationMeta{
			TotalRecords: total,
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			Limit:        int(pagination.Limit),
		},
	}, nil
}

func (s *ArtistService) GetArtistByID(ctx context.Context, id string) (dto.ArtistResponse, error) {
	artist, err := s.repo.GetArtistByID(ctx, id)
	if err != nil {
		return dto.ArtistResponse{}, err
	}

	return dto.ArtistResponse{
		ID:       artist.ID,
		Name:     artist.Name,
		Bio:      artist.Bio.String,
		CoverURL: artist.CoverUrl.String,
	}, nil
}

func (s *ArtistService) UpdateArtistByID(ctx context.Context, id string, req dto.UpdateArtistRequest) (dto.ArtistResponse, error) {
	artist, err := s.repo.UpdateArtistByID(ctx, repository.UpdateArtistByIDParams{
		Name:     req.Name,
		Bio:      req.Bio,
		CoverUrl: req.CoverURL,
		ID:       id,
	})
	if err != nil {
		return dto.ArtistResponse{}, err
	}

	return dto.ArtistResponse{
		ID:       artist.ID,
		Name:     artist.Name,
		Bio:      artist.Bio.String,
		CoverURL: artist.CoverUrl.String,
	}, nil
}

func (s *ArtistService) DeleteArtistByID(ctx context.Context, id string) (dto.ArtistResponse, error) {
	artist, err := s.repo.DeleteArtistByID(ctx, id)
	if err != nil {
		return dto.ArtistResponse{}, err
	}

	return dto.ArtistResponse{
		ID:       artist.ID,
		Name:     artist.Name,
		Bio:      artist.Bio.String,
		CoverURL: artist.CoverUrl.String,
	}, nil
}
