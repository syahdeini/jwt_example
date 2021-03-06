package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// if the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// For any other type of error, return a bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT String from the cookie
	tknStr := c.Value

	// Intialize a new instance of 'Claims'
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`
	// Note that we are passing the key in this method as well.
	// if the token is invalid (if it has expired according to the expiry time we set on sign in)
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s", claims.username)))
}
