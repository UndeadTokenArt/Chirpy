package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	corsMux := middlewareCors(router)
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// Mounted sub router for api access and Admin access
	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

	// handling the /app route
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	// handling the /admin routes
	adminRouter.Get("/metrics", apiCfg.metricsHtml)

	// handling the /api routes
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Get("/healthz", apiCfg.handlerHealthZ)
	apiRouter.Post("/chirps", apiCfg.handlerPostChirp)
	apiRouter.Get("/chirps", apiCfg.handlerGetChirp)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	fmt.Println("Server is running on http://localhost:8080")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

}
