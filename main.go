package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)
	corsMux := middlewareCors(router)

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	},
	)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	fmt.Println("Server is running on http://localhost:8080")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
