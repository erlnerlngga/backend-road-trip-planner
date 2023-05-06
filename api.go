package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type APIServer struct {
	listenAddr string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/register", makeHTTPHandleFunc(s.handleSignUp))

	log.Println("Server running in Port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "hello from the server"})
}

func (s *APIServer) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleVerifySignIn(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiError) handleLogout(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiError) handleGetDataDestination(w http.ResponseWriter, r *http.Request) error {
	// check if there is city in database or not

	// if yes get data from database

	// if not call data with go routines

	return nil
}

// Function Helper
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

// error handling
type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}

	}
}
