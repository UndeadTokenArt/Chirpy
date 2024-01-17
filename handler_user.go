package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogins(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters")
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:    user.ID,
		Email: user.Email,
	})
}

// HandleUserPost handles incoming HTTP requests for user registration
func (cfg *apiConfig) HandleUserPost(w http.ResponseWriter, r *http.Request) {
	// Define a struct to represent the expected parameters in the JSON request body.
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Timeout  int    `json:"expires_in_seconds"`
	}

	// Create a JSON decoder to read the request body.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	// Decode the JSON parameters into the 'params' struct.
	err := decoder.Decode(&params)
	if err != nil {
		// If decoding fails, respond with an internal server error.
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Convert the password to a hashed version using bcrypt.
	bytePass := []byte(params.Password)
	hashPass, err := bcrypt.GenerateFromPassword(bytePass, 1)
	if err != nil {
		// If password hashing fails, respond with an internal server error.
		respondWithError(w, http.StatusInternalServerError, "password did not hash")
		return
	}

	// Create a new user in the database with the hashed password.
	user, err := cfg.DB.CreateUser(params.Email, string(hashPass))
	if err != nil {
		// If user creation fails, respond with an internal server error.
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(params.Timeout))),
		Subject:   string(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(cfg.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token did not sign correctly")
	}

	// Respond with the newly created user details in JSON format.
	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
		Token: signedToken,
	})
}
