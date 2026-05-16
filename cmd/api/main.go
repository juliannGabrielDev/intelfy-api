package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/juliannGabrielDev/intelfy-api/internal/handler"
	"github.com/juliannGabrielDev/intelfy-api/internal/repository"
	"github.com/juliannGabrielDev/intelfy-api/internal/service"
)

func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Configure Postgres connection pool
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	// 3. Verify connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Postgres is not responding: %v\n", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	// 4. Initialize Repository
	queries := repository.New(dbPool)

	// 5. Initialize Services
	songService := service.NewSongService(queries)
	userService := service.NewUserService(queries)
	playlistService := service.NewPlaylistService(queries)
	genreService := service.NewGenreService(queries)
	albumService := service.NewAlbumService(queries)
	artistService := service.NewArtistService(queries)

	// 6. Initialize Handlers
	songHandler := handler.NewSongHandler(songService)
	userHandler := handler.NewUserHandler(userService)
	playlistHandler := handler.NewPlaylistHandler(playlistService)
	genreHandler := handler.NewGenreHandler(genreService)
	albumHandler := handler.NewAlbumHandler(albumService)
	artistHandler := handler.NewArtistHandler(artistService)

	// 7. Setup Router
	router := handler.NewRouter(songHandler, albumHandler, genreHandler, userHandler, playlistHandler, artistHandler)

	// 8. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
