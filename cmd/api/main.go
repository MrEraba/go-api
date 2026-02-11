package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ivan-almanza/notes-api/internal/api"
	"github.com/ivan-almanza/notes-api/internal/auth"
	"github.com/ivan-almanza/notes-api/internal/config"
	"github.com/ivan-almanza/notes-api/internal/store"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Configure Auth Secret
	auth.SetSecret(cfg.JWTSecret)

	// 3. Connect to Database
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 4. Initialize Store and Handlers
	postgresStore := store.NewPostgresStore(db)
	authHandler := api.NewAuthHandler(postgresStore)
	notesHandler := api.NewNotesHandler(postgresStore)

	// 5. Setup Router
	mux := http.NewServeMux()

	// Auth Routes
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// Notes Routes (Protected)
	mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
		api.WithAuth(http.HandlerFunc(notesHandler.CreateNote)).ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
		api.WithAuth(http.HandlerFunc(notesHandler.GetNotes)).ServeHTTP(w, r)
	})

	// 6. Start Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
