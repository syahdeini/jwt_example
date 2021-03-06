package main

import (
	// ...
	// Import the jwt-go library
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	// ...
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1":  "password1",
	"user2:": "password2",
}

// Create a struct to read the username and password from the request body
type Credential struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT
// We add jwt.StandartClaims as a embedeed type, to provide fields like expiry time
type Claims struct {
	username string `json: "username"`
	jwt.StandardClaims
}

// Create the Signin handler
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credential
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if not, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which include the username and expiry time
	claims := &Claims{
		username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// in JWT, the expirty time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT String
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// if there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
