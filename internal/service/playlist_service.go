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

type PlaylistService struct {
	repo *repository.Queries
}

func NewPlaylistService(repo *repository.Queries) *PlaylistService {
	return &PlaylistService{repo: repo}
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest) (*dto.PlaylistResponse, error) {
	if req.Name == "" {
		return nil, errors.New("playlist name can't be empty")
	}

	newID, err := nanoid.GenerateID()
	if err != nil {
		return nil, err
	}

	playlist, err := s.repo.CreatePlaylist(ctx, repository.CreatePlaylistParams{
		ID:   newID,
		Name: req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		UserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &dto.PlaylistResponse{
		ID:          playlist.ID,
		Name:        playlist.Name,
		Description: playlist.Description.String,
		UserID:      playlist.UserID,
		CreatedAt:   playlist.CreatedAt.Time,
	}, nil
}

func (s *PlaylistService) GetPlaylistsByUser(ctx context.Context, req dto.GetPlaylistsRequest) (*dto.PaginatedResponse[dto.PlaylistResponse], error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	playlists, err := s.repo.GetPlaylistsByUser(ctx, repository.GetPlaylistsByUserParams{
		UserID: req.UserID,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	resData := make([]dto.PlaylistResponse, len(playlists))
	for i, playlist := range playlists {
		resData[i] = s.mapToPlaylistResponse(playlist)
	}

	total, err := s.repo.CountPlaylistsByUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))
	currentPage := int(req.Offset/req.Limit) + 1

	return &dto.PaginatedResponse[dto.PlaylistResponse]{
		Data: resData,
		Meta: dto.PaginationMeta{
			TotalRecords: total,
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			Limit:        int(req.Limit),
		},
	}, nil
}

func (s *PlaylistService) GetPlaylistByID(ctx context.Context, id string) (*dto.GetPlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylistByID(ctx, id)
	if err != nil {
		return nil, err
	}

	songs, err := s.repo.GetSongsByPlaylistID(ctx, id)
	if err != nil {
		return nil, err
	}

	songResponses := make([]dto.SongResponse, len(songs))
	for i, song := range songs {
		songResponses[i] = dto.SongResponse{
			ID:       song.ID,
			Name:     song.Name,
			Duration: dto.SongDuration{Duration: time.Duration(song.DurationSeconds * float64(time.Second))},
			AudioURL: song.AudioUrl,
			AlbumID:  song.AlbumID,
			GenreID:  song.GenreID.String,
		}
	}

	return &dto.GetPlaylistResponse{
		Playlist: s.mapToPlaylistResponse(playlist),
		Songs:    songResponses,
	}, nil
}

func (s *PlaylistService) AddSongToPlaylist(ctx context.Context, req dto.AddSongToPlaylistRequest) (*dto.AddSongToPlaylistResponse, error) {
	err := s.repo.AddSongToPlaylist(ctx, repository.AddSongToPlaylistParams{
		PlaylistID: req.PlaylistID,
		SongID:     req.SongID,
	})
	if err != nil {
		return nil, err
	}

	playlist, err := s.repo.GetPlaylistByID(ctx, req.PlaylistID)
	if err != nil {
		return nil, err
	}

	song, err := s.repo.GetSongByID(ctx, req.SongID)
	if err != nil {
		return nil, err
	}

	totalSongs, err := s.repo.CountSongsInPlaylist(ctx, req.PlaylistID)
	if err != nil {
		return nil, err
	}

	return &dto.AddSongToPlaylistResponse{
		PlaylistID: playlist.ID,
		TotalSongs: int(totalSongs),
		AddedSong: dto.SongResponse{
			ID:       song.ID,
			Name:     song.Name,
			Duration: dto.SongDuration{Duration: time.Duration(song.DurationSeconds * float64(time.Second))},
			AudioURL: song.AudioUrl,
			AlbumID:  song.AlbumID,
			GenreID:  song.GenreID.String,
		},
	}, nil
}

func (s *PlaylistService) ClearPlaylist(ctx context.Context, playlistID string) error {
	err := s.repo.ClearPlaylist(ctx, playlistID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id string) (*dto.PlaylistResponse, error) {
	playlist, err := s.repo.DeletePlaylistByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PlaylistResponse{
		ID:          playlist.ID,
		Name:        playlist.Name,
		Description: playlist.Description.String,
		UserID:      playlist.UserID,
		CreatedAt:   playlist.CreatedAt.Time,
	}, nil
}

func (s *PlaylistService) mapToPlaylistResponse(playlist repository.Playlist) dto.PlaylistResponse {
	return dto.PlaylistResponse{
		ID:          playlist.ID,
		Name:        playlist.Name,
		Description: playlist.Description.String,
		UserID:      playlist.UserID,
		CreatedAt:   playlist.CreatedAt.Time,
	}
}
