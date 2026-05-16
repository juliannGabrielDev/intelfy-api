package service

import (
	"context"
	"errors"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juliannGabrielDev/intelfy-api/internal/dto"
	"github.com/juliannGabrielDev/intelfy-api/internal/repository"
	"github.com/juliannGabrielDev/intelfy-api/pkg/nanoid"
)

type AlbumService struct {
	repo *repository.Queries
}

func NewAlbumService(repo *repository.Queries) *AlbumService {
	return &AlbumService{repo: repo}
}

func (s *AlbumService) CreateAlbum(ctx context.Context, req dto.CreateAlbumRequest) (dto.AlbumResponse, error) {
	newID, err := nanoid.GenerateID()
	if err != nil {
		return dto.AlbumResponse{}, err
	}

	album, err := s.repo.CreateAlbum(ctx, repository.CreateAlbumParams{
		ID:   newID,
		Name: req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		CoverUrl: pgtype.Text{
			String: req.CoverURL,
			Valid:  req.CoverURL != "",
		},
		ArtistID: req.ArtistID,
		ReleaseDate: pgtype.Date{
			Time:  req.ReleaseDate.Time,
			Valid: !req.ReleaseDate.Time.IsZero(),
		},
	})

	if err != nil {
		return dto.AlbumResponse{}, err
	}

	return s.mapToAlbumResponse(album), nil
}

func (s *AlbumService) GetAlbums(ctx context.Context, filters dto.AlbumFilters) (*dto.PaginatedResponse[dto.AlbumResponse], error) {
	if filters.Limit <= 0 {
		filters.Limit = 10
	}

	total, err := s.repo.CountAlbums(ctx, filters.ArtistID)

	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dto.PaginatedResponse[dto.AlbumResponse]{
			Data: []dto.AlbumResponse{},
			Meta: dto.PaginationMeta{
				TotalRecords: 0,
				CurrentPage:  1,
				TotalPages:   0,
				Limit:        int(filters.Limit),
			},
		}, nil
	}

	albums, err := s.repo.GetAlbums(ctx, repository.GetAlbumsParams{
		ArtistID: filters.ArtistID,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	})

	if err != nil {
		return nil, err
	}

	resData := make([]dto.AlbumResponse, len(albums))
	for i, album := range albums {
		resData[i] = s.mapToAlbumResponse(repository.Album(album))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filters.Limit)))
	currentPage := int(filters.Offset/filters.Limit) + 1

	return &dto.PaginatedResponse[dto.AlbumResponse]{
		Data: resData,
		Meta: dto.PaginationMeta{
			TotalRecords: total,
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			Limit:        int(filters.Limit),
		},
	}, nil
}

func (s *AlbumService) GetAlbumByID(ctx context.Context, id string) (dto.AlbumResponse, error) {
	album, err := s.repo.GetAlbumByID(ctx, id)
	if err != nil {
		return dto.AlbumResponse{}, err
	}

	return s.mapToAlbumResponse(album), nil
}

func (s *AlbumService) UpdateAlbumByID(ctx context.Context, id string, req dto.UpdateAlbumRequest) error {
	if req.Name != nil && *req.Name == "" {
		return errors.New("album name cannot be empty")
	}

	if req.ArtistID != nil && *req.ArtistID == "" {
		return errors.New("artist id cannot be empty")
	}

	var releaseDate pgtype.Date
	if req.ReleaseDate != nil {
		releaseDate = pgtype.Date{
			Time:  req.ReleaseDate.Time,
			Valid: true,
		}
	}

	return s.repo.UpdateAlbumByID(ctx, repository.UpdateAlbumByIDParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		CoverUrl:    req.CoverURL,
		ArtistID:    req.ArtistID,
		ReleaseDate: releaseDate,
	})
}

func (s *AlbumService) DeleteAlbumByID(ctx context.Context, id string) error {
	return s.repo.DeleteAlbumByID(ctx, id)
}

func (s *AlbumService) mapToAlbumResponse(album repository.Album) dto.AlbumResponse {
	return dto.AlbumResponse{
		ID:          album.ID,
		Name:        album.Name,
		Description: album.Description.String,
		CoverURL:    album.CoverUrl.String,
		ArtistID:    album.ArtistID,
		ReleaseDate: album.ReleaseDate.Time,
	}
}
