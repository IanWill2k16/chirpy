package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/IanWill2k16/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	envPlatform    string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	envPlatformVar := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("error connect to database: %v", err)
	}

	apiCfg := &apiConfig{
		dbQueries:   database.New(db),
		envPlatform: envPlatformVar,
	}

	const port = "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app")))))

	mux.HandleFunc("GET /api/healthz", healthResp)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/login", apiCfg.loginUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getSingleChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.countResp)
	mux.HandleFunc("POST /admin/reset", apiCfg.adminReset)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
