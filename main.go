package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	corsMux := middlewareCors(router)

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	// Mounted sub router for api access and Admin access
	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)
	// handling the app route
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)
	//
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
