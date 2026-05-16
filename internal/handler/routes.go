package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/juliannGabrielDev/intelfy-api/internal/middleware"
)

func NewRouter(songH *SongHandler, albumH *AlbumHandler, genreH *GenreHandler, userH *UserHandler, playlistH *PlaylistHandler, artistH *ArtistHandler) *chi.Mux {
	r := chi.NewRouter()

	// Standard middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(customMiddleware.RateLimiter)

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Songs
		r.Route("/songs", func(r chi.Router) {
			r.Post("/", songH.CreateSong)
			r.Get("/", songH.GetSongs)
			r.Get("/{id}", songH.GetSongByID)
			r.Patch("/{id}", songH.UpdateSongByID)
			r.Delete("/{id}", songH.DeleteSongByID)
		})

		// Albums
		r.Route("/albums", func(r chi.Router) {
			r.Post("/", albumH.CreateAlbum)
			r.Get("/", albumH.GetAlbums)
			r.Get("/{id}", albumH.GetAlbumByID)
			r.Patch("/{id}", albumH.UpdateAlbumByID)
			r.Delete("/{id}", albumH.DeleteAlbumByID)
		})

		// Genres
		r.Route("/genres", func(r chi.Router) {
			r.Post("/", genreH.CreateGenre)
			r.Get("/", genreH.GetGenres)
			r.Get("/{id}", genreH.GetGenreByID)
			r.Patch("/{id}", genreH.UpdateGenreByID)
			r.Delete("/{id}", genreH.DeleteGenreByID)
		})

		// Artists
		r.Route("/artists", func(r chi.Router) {
			r.Post("/", artistH.CreateArtist)
			r.Get("/", artistH.GetArtists)
			r.Get("/{id}", artistH.GetArtistByID)
			r.Patch("/{id}", artistH.UpdateArtistByID)
			r.Delete("/{id}", artistH.DeleteArtistByID)
		})

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", userH.Register)
			r.Post("/login", userH.Login)

			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RequireAuth)
				r.Get("/", userH.GetUsers)
				r.Get("/{id}", userH.GetUserByID)
				r.Patch("/{id}", userH.UpdateUser)
			})
		})

		// Playlists
		r.Route("/playlists", func(r chi.Router) {
			r.Post("/", playlistH.CreatePlaylist)
			r.Get("/users/{user_id}", playlistH.GetPlaylistsByUser)
			r.Get("/{id}", playlistH.GetPlaylistByID)
			r.Post("/songs", playlistH.AddSongToPlaylist)
			r.Delete("/{id}/songs", playlistH.ClearPlaylist)
			r.Delete("/{id}", playlistH.DeletePlaylist)
		})
	})

	return r
}
