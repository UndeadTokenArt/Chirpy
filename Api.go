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

func (cfg *apiConfig) handlerHealthZ(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	// returns the current list of chirps
}

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	// expected structure of request
	type paramaters struct {
		Body string `json:"body"`
	}

	// structure of the response
	type returnVals struct {
		Id   int    `json:"id"`
		Body string `json:"Body"`
	}
	// create a response interface
	respBody := returnVals{}

	// json Decoding so we can mess with the data
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding paramaters: %v", err)
		w.WriteHeader(500)
	}

	// If the length of the Chirp is too long
	if len(params.Body) > 140 {
		w.WriteHeader(400)
	} else

	// fill in the fields and prep for response
	{
		cleaned := censorChirps(params.Body)
		respBody.Body = cleaned
		respBody.Id = assignID(respBody, int)
		w.WriteHeader(201)
	}
	respondWithJson(w, 201, respBody)
}

func censorChirps(body string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func respondWithJson(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(status)
	w.Write(dat)

}

func assignID(r interface{}) int {

	// Gets the current highest ID number from the database
}
