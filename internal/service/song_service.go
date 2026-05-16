package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juliannGabrielDev/intelfy-api/internal/dto"
	"github.com/juliannGabrielDev/intelfy-api/internal/repository"
	"github.com/juliannGabrielDev/intelfy-api/pkg/nanoid"
)

type SongService struct {
	repo *repository.Queries
}

func NewSongService(repo *repository.Queries) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) CreateSong(ctx context.Context, req dto.CreateSongRequest) (dto.SongResponse, error) {
	newID, err := nanoid.GenerateID()
	if err != nil {
		return dto.SongResponse{}, err
	}

	song, err := s.repo.CreateSong(ctx, repository.CreateSongParams{
		ID:              newID,
		Name:            req.Name,
		DurationSeconds: req.Duration.Seconds(),
		AudioUrl:        req.AudioURL,
		AlbumID:         req.AlbumID,
		GenreID: pgtype.Text{
			String: req.GenreID,
			Valid:  req.GenreID != "",
		},
	})

	if err != nil {
		return dto.SongResponse{}, err
	}

	return s.mapToSongResponse(song), nil
}

func (s *SongService) GetSongs(ctx context.Context, filters dto.SongFilters) (*dto.PaginatedResponse[dto.SongResponse], error) {
	if filters.Limit <= 0 {
		filters.Limit = 10
	}

	genreID := pgtype.Text{
		String: filters.GenreID,
		Valid:  filters.GenreID != "",
	}

	total, err := s.repo.CountSongs(ctx, repository.CountSongsParams{
		AlbumID:  filters.AlbumID,
		ArtistID: filters.ArtistID,
		GenreID:  genreID,
	})
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dto.PaginatedResponse[dto.SongResponse]{
			Data: []dto.SongResponse{},
			Meta: dto.PaginationMeta{
				TotalRecords: 0,
				CurrentPage:  1,
				TotalPages:   1,
				Limit:        int(filters.Limit),
			},
		}, nil
	}

	songs, err := s.repo.GetSongs(ctx, repository.GetSongsParams{
		AlbumID:  filters.AlbumID,
		ArtistID: filters.ArtistID,
		GenreID:  genreID,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	})

	if err != nil {
		return nil, err
	}

	resData := make([]dto.SongResponse, len(songs))
	for i, song := range songs {
		resData[i] = s.mapToSongResponse(repository.Song(song))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filters.Limit)))
	currentPage := int(filters.Offset/filters.Limit) + 1

	return &dto.PaginatedResponse[dto.SongResponse]{
		Data: resData,
		Meta: dto.PaginationMeta{
			TotalRecords: total,
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			Limit:        int(filters.Limit),
		},
	}, nil
}

func (s *SongService) GetSongByID(ctx context.Context, id string) (dto.SongResponse, error) {
	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return dto.SongResponse{}, err
	}

	return s.mapToSongResponse(song), nil
}

func (s *SongService) UpdateSongByID(ctx context.Context, id string, req dto.UpdateSongRequest) error {
	// Verify song exists
	_, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return errors.New("song not found")
	}

	if req.Name != nil && *req.Name == "" {
		return errors.New("song name cannot be empty")
	}

	if req.Duration != nil && req.Duration.Seconds() <= 0 {
		return errors.New("duration must be greater than zero")
	}

	if req.AudioURL != nil && *req.AudioURL == "" {
		return errors.New("audio url cannot be empty")
	}

	if req.AlbumID != nil && *req.AlbumID == "" {
		return errors.New("album id cannot be empty")
	}

	return s.repo.UpdateSongByID(ctx, repository.UpdateSongByIDParams{
		ID:              id,
		Name:            req.Name,
		DurationSeconds: durationPointerToSeconds(req.Duration),
		AudioUrl:        req.AudioURL,
		AlbumID:         req.AlbumID,
		GenreID:         req.GenreID,
	})
}

func (s *SongService) DeleteSongByID(ctx context.Context, id string) error {
	// Verify song exists
	_, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return errors.New("song not found")
	}

	return s.repo.DeleteSongByID(ctx, id)
}

func (s *SongService) mapToSongResponse(song repository.Song) dto.SongResponse {
	return dto.SongResponse{
		ID:       song.ID,
		Name:     song.Name,
		Duration: dto.SongDuration{Duration: time.Duration(song.DurationSeconds * float64(time.Second))},
		AudioURL: song.AudioUrl,
		AlbumID:  song.AlbumID,
		GenreID:  song.GenreID.String,
	}
}

func durationPointerToSeconds(d *dto.SongDuration) interface{} {
	if d == nil {
		return nil
	}

	return d.Seconds()
}
