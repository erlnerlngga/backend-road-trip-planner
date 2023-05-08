package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	// "strconv"

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

	router.Post("/signup", makeHTTPHandleFunc(s.handleSignUp))
	router.Post("/signin", makeHTTPHandleFunc(s.handleSignIn))
	router.Get("/signin/{token}", makeHTTPHandleFunc(s.handleVerifySignIn))
	router.Get("/logout", makeHTTPHandleFunc(s.handleLogout))

	router.Get("/destination/{city}", makeHTTPHandleFunc(s.handleGetAllDestination))
	router.Get("/destination/{city}-{destination_id}", makeHTTPHandleFunc(s.handleGetDestination))

	// create new bookmark
	router.Post("/bookmark", makeHTTPHandleFunc(s.handleCreateNewBookmark))

	// save data into bookmark
	router.Post("/bookmark/save", makeHTTPHandleFunc(s.handleSaveIntoBookmark))

	// create new bookmark and save
	router.Post("/bookmark/create-and-save", makeHTTPHandleFunc(s.handleCreateAndSaveIntoBookmark))

	// get all bookmark name
	router.Get("/bookmark/{user_id}", makeHTTPHandleFunc(s.handleGetBookmarkName))

	// handle get all data from bookmark call save_user table base on bookmark_id and join with destination table
	router.Get("/bookmark/{bookmark_id}", makeHTTPHandleFunc(s.handleGetBookmarkData))

	// handle updated bookmark name
	router.Put("/boomark/{bookmark_id}", makeHTTPHandleFunc(s.handleBookmarkUpdateName))

	// handle delete bookmark name 
	router.Delete("/bookmark/{bookmark_id}", makeHTTPHandleFunc(s.handleDeleteBookmarkName))

	// handle delete data bookmark
	router.Delete("/bookmark/{destination_book_id}", makeHTTPHandleFunc(s.handleDeleteBookmarkDestination))

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

func (s *APIServer) handleLogout(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle GET ALL DATA DESTINATION
func (s *APIServer) handleGetAllDestination(w http.ResponseWriter, r *http.Request) error {
	// get param city
	// cityParam := chi.URLParam(r, "city")

	// check if there is city in database or not

	// if yes get data from database

	// if not call data with go routines

	return nil
}

// handle get ONE DESTINATION with BUNCH of IMAGE
func (s *APIServer) handleGetDestination(w http.ResponseWriter, r *http.Request) error {
	// get param city and destination_id
	// cityParam := chi.URLParam(r, "city")
	// destination_idParam := chi.URLParam(r, "destination_id")

	// call destination table to get name and url

	// and call image table to get all of image

	return nil
}

// handle create new bookmark
func (s *APIServer) handleCreateNewBookmark(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle get all bookmark name
func (s *APIServer) handleGetBookmarkName(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle save data into bookmark
func (s *APIServer) handleSaveIntoBookmark(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle create new bookmark and save data
func (s *APIServer) handleCreateAndSaveIntoBookmark(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle update bookmark name
func (s *APIServer) handleBookmarkUpdateName(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle get bookmark data
func (s *APIServer) handleGetBookmarkData(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle delete bookmark name
func (s *APIServer) handleDeleteBookmarkName(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handle delete bookmark destination
func (s *APIServer) handleDeleteBookmarkDestination(w http.ResponseWriter, r *http.Request) error {
	// bookmark_id, err := strconv.Atoi(chi.URLParam(r, "bookmark_id"))

	// if err != nil {
	// 	return err
	// }

	// destination_id, err :=	strconv.Atoi(chi.URLParam(r, "destination_book_id"))

	// if err != nil {
	// 	return err
	// }

	// createUser_Save := &CreateNewUser_SaveType{
	// 	Destination_ID: destination_id,
	// 	Bookmark_ID: bookmark_id,
	// }

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
