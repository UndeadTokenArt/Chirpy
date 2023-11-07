package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
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

func (cfg *apiConfig) metricsHtml(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	templateHtml := `
	<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	
	</html>
	`
	w.Write([]byte(fmt.Sprintf(templateHtml, cfg.fileserverHits)))
}

func (cfg *apiConfig) handleHealthZ(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handleValidate(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")

	type paramaters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Body  string `json:"body"`
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}
	respBody := returnVals{}

	decoder := json.NewDecoder(r.Body)
	params := paramaters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding paramaters: %v", err)
		w.WriteHeader(500)
	}

	if len(params.Body) > 140 {
		respBody.Error = "chirp is too long"
		respBody.Valid = false
		w.WriteHeader(400)
	} else {
		cleaned := censorChirps(params.Body)
		respBody.Body = cleaned
		respBody.Valid = true
		w.WriteHeader(200)
	}

	// form the response back to the client from the respBody
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)

}
func censorChirps(chirp string) string {
	var badwords []string
	badwords = append(badwords, "kerfuffle", "sharbert", "fornax")
	for _, badword := range badwords {
		fmt.Printf("bad word: %s and current chirp: %s\n", badword, chirp)
		chirp = strings.Replace(chirp, badword, "****", -1)
	}
	return chirp
}
