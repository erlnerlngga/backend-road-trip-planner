package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// create JWT
func CreateJWT(user_id int) (string, error) {
	// declare expiration time with 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// declare jwt claims
	claims := &ClaimsType{
		User_ID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// declare token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// create jwt string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// MIDDLEWARE TO HANDLE JWT VERIFICATION
func WithJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		// get auth header
		authHeader := r.Header.Get("Authorization")

		// sanity check
		if authHeader == "" {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "no auth header"})
			return 
		} 

		// split the header space
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "invalid auth header"})
			return 
		}

		tokenString := headerParts[1]

		// init claims
		claims := new(ClaimsType)

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Signature Invalid"})
				return
			}

			if strings.HasPrefix(err.Error(), "token is expired by") {
				WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "expired token"})
				return
			}

			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "token invalid"})
			return
		}

		next.ServeHTTP(w, r)
	})
}