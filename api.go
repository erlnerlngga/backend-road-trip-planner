package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type APIServer struct {
	listenAddr string
	store      Storage
	user_id    string
}

func NewApiServer(listenAddr string, storage Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      storage,
	}
}

func (s *APIServer) Run() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://roadtrip.vercel.app/", "https://roadtrip-laannen-gmailcom.vercel.app", "https://roadtrip-q7ki6cz9s-laannen-gmailcom.vercel.app", "https://roadtrip-git-main-laannen-gmailcom.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	router.Get("/", makeHTTPHandleFunc(s.handleWelcome))
	router.Post("/signup", makeHTTPHandleFunc(s.handleSignUp))
	router.Post("/signin", makeHTTPHandleFunc(s.handleSignIn))
	router.Get("/auth/{token}", makeHTTPHandleFunc(s.handleVerifySignIn))

	router.Group(func(r chi.Router) {
		r.Use(WithJWTAuth)
		r.Get("/logout", makeHTTPHandleFunc(s.handleLogout))
		r.Get("/destination/{city}", makeHTTPHandleFunc(s.handleGetAllDestination))
		r.Get("/destination/specific/{destination_id}", makeHTTPHandleFunc(s.handleGetDestination))
		r.Post("/bookmark", makeHTTPHandleFunc(s.handleCreateNewBookmark))
		r.Post("/bookmark/save", makeHTTPHandleFunc(s.handleSaveIntoBookmark))
		r.Post("/bookmark/create-and-save", makeHTTPHandleFunc(s.handleCreateAndSaveIntoBookmark))
		r.Get("/bookmark", makeHTTPHandleFunc(s.handleGetBookmarkName))
		r.Get("/bookmark/specific/{bookmark_id}", makeHTTPHandleFunc(s.handleGetBookmarkData))
		r.Put("/bookmark/{bookmark_id}", makeHTTPHandleFunc(s.handleBookmarkUpdateName))
		r.Delete("/bookmark/{bookmark_id}", makeHTTPHandleFunc(s.handleDeleteBookmarkName))
		r.Delete("/bookmark/specific/{destination_book_id}", makeHTTPHandleFunc(s.handleDeleteBookmarkDestination))
	})

	log.Println("Server running in Port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleWelcome(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "welcome to roadtrip"})
}

func (s *APIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	newAccount := new(SignUpType)
	if err := json.NewDecoder(r.Body).Decode(newAccount); err != nil {
		return err
	}

	defer r.Body.Close()

	account, err := s.store.SignUp(newAccount)
	if err != nil {
		return err
	}

	tokenStr, err := CreateJWT(account.User_ID)
	if err != nil {
		return err
	}

	if err := SendMAIL(account.Email, account.User_Name, tokenStr); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (s *APIServer) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	email := new(SignInType)
	if err := json.NewDecoder(r.Body).Decode(email); err != nil {
		log.Println("1. handleSignIn", err)
		return err
	}
	defer r.Body.Close()

	account, err := s.store.CheckEmail(email.Email)
	if err != nil {
		log.Println("2. handleSignIn", err)
		return err
	}

	// create token
	tokenStr, err := CreateJWT(account.User_ID)
	if err != nil {
		log.Println("3. handleSignIn", err)
		return err
	}

	if err := SendMAIL(account.Email, account.User_Name, tokenStr); err != nil {
		log.Println("4. handleSignIn", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (s *APIServer) handleVerifySignIn(w http.ResponseWriter, r *http.Request) error {
	tokenStr := chi.URLParam(r, "token")

	claims := new(ClaimsType)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("1. handleVerifySignIn", err)
			return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Signature Invalid"})
		}

		log.Println("2. handleVerifySignIn", err)
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: err.Error()})
	}

	if !token.Valid {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "token invalid"})
	}

	s.user_id = claims.User_ID
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "token": tokenStr})
}

func (s *APIServer) handleLogout(w http.ResponseWriter, r *http.Request) error {
	s.user_id = ""
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "Logout success"})
}

// handle GET ALL DATA DESTINATION
func (s *APIServer) handleGetAllDestination(w http.ResponseWriter, r *http.Request) error {
	// get param city
	param := chi.URLParam(r, "city")
	cityParam, err := url.QueryUnescape(param)
	if err != nil {
		log.Println("1. handleGetAllDestination", err)
		return err
	}

	// check if there is city in database or not
	city, err := s.store.CheckCity(cityParam)
	if err != nil {
		log.Println("2. handleGetAllDestination", err)
		return err
	}

	// get list destination base on city
	allDestination, err := s.store.GetAllDestination(city.City_ID)
	if err != nil {
		log.Println("3. handleGetAllDestination", err)
		return err
	}

	sendAllData := &SendAllDestinationType{
		City_Name:        city.City_Name,
		City_Lat:         city.City_Lat,
		City_Long:        city.City_Long,
		List_Destination: allDestination,
	}

	return WriteJSON(w, http.StatusOK, sendAllData)
}

// handle get ONE DESTINATION with BUNCH of IMAGE
func (s *APIServer) handleGetDestination(w http.ResponseWriter, r *http.Request) error {
	// get param  and destination_id
	destination_idParam := chi.URLParam(r, "destination_id")

	// call destination table to get name and url
	destination, err := s.store.GetDestination(destination_idParam)
	if err != nil {
		log.Println("1. handleGetDestination", err)
		return err
	}

	// and call image table to get all of image
	images, err := s.store.GetAllImages(destination.Destination_ID)
	if err != nil {
		log.Println("2. handleGetDestination", err)
		return err
	}

	sendData := &SendSpecificDestinationType{
		Destination_ID:   destination.Destination_ID,
		Destination_Name: destination.Destination_Name,
		Destination_URL:  destination.Destination_URL,
		List_Image:       images,
	}

	return WriteJSON(w, http.StatusOK, sendData)
}

// handle create new bookmark
func (s *APIServer) handleCreateNewBookmark(w http.ResponseWriter, r *http.Request) error {
	// read data from the body
	book := new(NewBookmarkType)
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		log.Println("1. handleCreateNewBookmark", err)
		return err
	}

	book.User_ID = s.user_id

	defer r.Body.Close()

	_, err := s.store.CreateNewBookmark(book)
	if err != nil {
		log.Println("2. handleCreateNewBookmark", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle get all bookmark name
func (s *APIServer) handleGetBookmarkName(w http.ResponseWriter, r *http.Request) error {
	bookmarks, err := s.store.GetAllBookmark(s.user_id)
	if err != nil {
		log.Println("3. handleCreateNewBookmark", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, bookmarks)
}

// handle save data into bookmark
func (s *APIServer) handleSaveIntoBookmark(w http.ResponseWriter, r *http.Request) error {
	newSave := new(CreateNewUser_SaveType)
	if err := json.NewDecoder(r.Body).Decode(newSave); err != nil {
		log.Println("1. handleSaveIntoBookmark", err)
		return err
	}

	if err := s.store.SaveBookmarkData(newSave); err != nil {
		log.Println("2. handleSaveIntoBookmark", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle create new bookmark and save data
func (s *APIServer) handleCreateAndSaveIntoBookmark(w http.ResponseWriter, r *http.Request) error {
	newBookReq := new(CreateBookmarkAndSaveType)
	if err := json.NewDecoder(r.Body).Decode(newBookReq); err != nil {
		log.Println("1. handleCreateAndSaveIntoBookmark", err)
		return err
	}

	// create bookmark
	newBookData := &NewBookmarkType{
		User_ID:       newBookReq.User_ID,
		Bookmark_Name: newBookReq.Bookmark_Name,
	}

	newBook, err := s.store.CreateNewBookmark(newBookData)
	if err != nil {
		log.Println("2. handleCreateAndSaveIntoBookmark", err)
		return err
	}

	// save data
	newSaveData := &CreateNewUser_SaveType{
		Destination_ID: newBookReq.Destination_ID,
		Bookmark_ID:    newBook.Bookmark_ID,
	}

	if err := s.store.SaveBookmarkData(newSaveData); err != nil {
		log.Println("3. handleCreateAndSaveIntoBookmark", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle update bookmark name
func (s *APIServer) handleBookmarkUpdateName(w http.ResponseWriter, r *http.Request) error {
	bookID := chi.URLParam(r, "bookmark_id")

	bookNewName := new(UpdateBookmarkNameType)
	if err := json.NewDecoder(r.Body).Decode(bookNewName); err != nil {
		log.Println("1. handleBookmarkUpdateName", err)
		return err
	}

	// update
	if err := s.store.UpdateBookmarkName(bookID, bookNewName); err != nil {
		log.Println("2. handleBookmarkUpdateName", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle get bookmark data
func (s *APIServer) handleGetBookmarkData(w http.ResponseWriter, r *http.Request) error {
	bookmark_id := chi.URLParam(r, "bookmark_id")

	user_save_data, err := s.store.GetAllDataByBookmark(bookmark_id)
	if err != nil {
		log.Println("1. handleGetBookmarkData", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, user_save_data)
}

// handle delete bookmark name
func (s *APIServer) handleDeleteBookmarkName(w http.ResponseWriter, r *http.Request) error {
	bookmark_id := chi.URLParam(r, "bookmark_id")

	if err := s.store.DeleteBookmark(bookmark_id); err != nil {
		log.Println("1. handleDeleteBookmarkName", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle delete bookmark destination
func (s *APIServer) handleDeleteBookmarkDestination(w http.ResponseWriter, r *http.Request) error {
	user_dave_id := chi.URLParam(r, "destination_book_id")

	if err := s.store.DeleteBookmarkData(user_dave_id); err != nil {
		log.Println("1. handleDeleteBookmarkDestination", err)
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
