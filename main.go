package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/undeadtokenart/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	SecretKey      []byte
}

func main() {
	godotenv.Load()
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	// sets the --debug flag for deleting the json "database" at run
	dbg := flag.Bool("debug", false, "enable bug mode")
	flag.Parse()

	if *dbg {
		os.Remove("database.json")
	}

	// sets the root file path and the port number to use
	const filepathRoot = "."
	const port = "8080"

	// if there is not a database.json file for holding the account and chirps, then it creates one.
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	// resetting the how many times the fileserver has been accessed
	// information needed for defining the database to use
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		SecretKey:      jwtSecret,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	// all the routes accessed from the api subroute
	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	apiRouter.Post("/login", apiCfg.handlerLogins)
	apiRouter.Post("/users", apiCfg.HandleUserPost)

	// sets the admin path and the subroutes of the admin path
	adminRouter := chi.NewRouter()
	router.Mount("/admin", adminRouter)
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	//
	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
