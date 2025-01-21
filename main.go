package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	apiCfg := &apiConfig{}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app")))))

	mux.HandleFunc("GET /api/healthz", healthResp)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.countResp)
	mux.HandleFunc("POST /admin/reset", apiCfg.countReset)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
