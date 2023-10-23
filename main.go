package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	mux.Handle("/", http.FileServer(http.Dir(".")))
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
