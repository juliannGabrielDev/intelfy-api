package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/juliannGabrielDev/intelfy-api/internal/middleware"
)

func NewRouter(songH *SongHandler, albumH *AlbumHandler, genreH *GenreHandler, userH *UserHandler, playlistH *PlaylistHandler, artistH *ArtistHandler, wsH *WSHandler) *chi.Mux {
	r := chi.NewRouter()

	// Standard middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(customMiddleware.RateLimiter)

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// WebSockets
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.RequireAuth)
			r.Get("/ws", wsH.HandleWS)
		})

		// Songs
		r.Route("/songs", func(r chi.Router) {
			r.Get("/", songH.GetSongs)
			r.Get("/{id}", songH.GetSongByID)

			// Artist only routes
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RequireAuth)
				r.Use(customMiddleware.RequireRole("artist"))
				r.Post("/", songH.CreateSong)
				r.Patch("/{id}", songH.UpdateSongByID)
				r.Delete("/{id}", songH.DeleteSongByID)
			})
		})

		// Albums
		r.Route("/albums", func(r chi.Router) {
			r.Get("/", albumH.GetAlbums)
			r.Get("/{id}", albumH.GetAlbumByID)

			// Artist only routes
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RequireAuth)
				r.Use(customMiddleware.RequireRole("artist"))
				r.Post("/", albumH.CreateAlbum)
				r.Patch("/{id}", albumH.UpdateAlbumByID)
				r.Delete("/{id}", albumH.DeleteAlbumByID)
			})
		})

		// Genres
		r.Route("/genres", func(r chi.Router) {
			r.Get("/", genreH.GetGenres)
			r.Get("/{id}", genreH.GetGenreByID)

			// Admin/Artist only? Let's keep it simple for now, maybe just Auth
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RequireAuth)
				r.Post("/", genreH.CreateGenre)
				r.Patch("/{id}", genreH.UpdateGenreByID)
				r.Delete("/{id}", genreH.DeleteGenreByID)
			})
		})

		// Artists
		r.Route("/artists", func(r chi.Router) {
			r.Get("/", artistH.GetArtists)
			r.Get("/{id}", artistH.GetArtistByID)

			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RequireAuth)
				r.Post("/", artistH.CreateArtist)
				r.Patch("/{id}", artistH.UpdateArtistByID)
				r.Delete("/{id}", artistH.DeleteArtistByID)
			})
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
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RequireAuth)
				r.Post("/", playlistH.CreatePlaylist)
				r.Get("/users/{user_id}", playlistH.GetPlaylistsByUser)
				r.Get("/{id}", playlistH.GetPlaylistByID)
				r.Post("/songs", playlistH.AddSongToPlaylist)
				r.Delete("/{id}/songs", playlistH.ClearPlaylist)
				r.Delete("/{id}", playlistH.DeletePlaylist)
			})
		})
	})

	// Static files (Protected)
	workDir, _ := os.Getwd()
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	filesDir := http.Dir(filepath.Join(workDir, uploadDir))

	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.RequireAuth)
		fileServer(r, "/"+uploadDir, filesDir)
	})

	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
