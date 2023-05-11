package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"time"

	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type APIServer struct {
	listenAddr string
	store      Storage
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
		AllowedOrigins:   []string{"http://127.0.0.1", "http://0.0.0.0"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	router.Post("/signup", makeHTTPHandleFunc(s.handleSignUp))
	router.Post("/signin", makeHTTPHandleFunc(s.handleSignIn))
	router.Get("/signin/{token}", makeHTTPHandleFunc(s.handleVerifySignIn))

	router.Get("/logout", makeHTTPHandleFunc(s.handleLogout))

	router.Get("/destination/{city}", makeHTTPHandleFunc(s.handleGetAllDestination))
	router.Get("/destination/specific/{destination_id}", makeHTTPHandleFunc(s.handleGetDestination))

	// create new bookmark
	router.Post("/bookmark", makeHTTPHandleFunc(s.handleCreateNewBookmark))

	// save data into bookmark
	router.Post("/bookmark/save", makeHTTPHandleFunc(s.handleSaveIntoBookmark))

	// create new bookmark and save
	router.Post("/bookmark/create-and-save", makeHTTPHandleFunc(s.handleCreateAndSaveIntoBookmark))

	// get all bookmark name
	router.Get("/bookmark/{user_id}", makeHTTPHandleFunc(s.handleGetBookmarkName))

	// handle get all data from bookmark call save_user table base on bookmark_id and join with destination table
	router.Get("/bookmark/specific/{bookmark_id}", makeHTTPHandleFunc(s.handleGetBookmarkData))

	// handle updated bookmark name
	router.Put("/bookmark/{bookmark_id}", makeHTTPHandleFunc(s.handleBookmarkUpdateName))

	// handle delete bookmark name
	router.Delete("/bookmark/{bookmark_id}", makeHTTPHandleFunc(s.handleDeleteBookmarkName))

	// handle delete data bookmark
	router.Delete("/bookmark/specific/{destination_book_id}", makeHTTPHandleFunc(s.handleDeleteBookmarkDestination))

	log.Println("Server running in Port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
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

	tokenStr, err := CreateJWT(account.ID)
	if err != nil {
		return err
	}

	if err := sendEmail(account.Email, tokenStr); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	email := new(SignInType)
	if err := json.NewDecoder(r.Body).Decode(email); err != nil {
		return err
	}
	defer r.Body.Close()

	account, err := s.store.CheckEmail(email.Email)
	if err != nil {
		return err
	}

	// create token
	tokenStr, err := CreateJWT(account.ID)
	if err != nil {
		return err
	}

	if err := sendEmail(account.Email, tokenStr); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleVerifySignIn(w http.ResponseWriter, r *http.Request) error {
	tokenStr := chi.URLParam(r, "token")

	claims := new(ClaimsType)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return WriteJSON(w, http.StatusOK, ApiError{Error: "Signature Invalid"})
		}
		return WriteJSON(w, http.StatusOK, ApiError{Error: err.Error()})
	}

	if !token.Valid {
		return WriteJSON(w, http.StatusOK, ApiError{Error: "token invalid"})
		//return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "token invalid"})
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: claims.ExpiresAt.Time,
		Domain:  "http://localhost:8080",
		Path:    "/",
	})
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "user_id": strconv.Itoa(claims.User_ID)})
}

func (s *APIServer) handleLogout(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
		MaxAge:  -1,
		Domain:  "http://localhost:8080",
		Path:    "/",
	})
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "Logout success"})
}

// handle GET ALL DATA DESTINATION
func (s *APIServer) handleGetAllDestination(w http.ResponseWriter, r *http.Request) error {
	// get param city
	param := chi.URLParam(r, "city")
	cityParam, err := url.QueryUnescape(param)
	if err != nil {
		return err
	}

	// check if there is city in database or not
	city, err := s.store.CheckCity(cityParam)
	if err != nil {
		return err
	}

	// get list destination base on city
	allDestination, err := s.store.GetAllDestination(city.ID)
	if err != nil {
		return err
	}

	sendAllData := &SendAllDestinationType{
		Name_City:        city.Name_City,
		Lat_City:         city.Lat_City,
		Long_City:        city.Long_City,
		List_Destination: allDestination,
	}

	return WriteJSON(w, http.StatusOK, sendAllData)
}

// handle get ONE DESTINATION with BUNCH of IMAGE
func (s *APIServer) handleGetDestination(w http.ResponseWriter, r *http.Request) error {
	// get param  and destination_id
	destination_idParam, err := strconv.Atoi(chi.URLParam(r, "destination_id"))
	if err != nil {
		return err
	}

	// call destination table to get name and url
	destination, err := s.store.GetDestination(destination_idParam)
	if err != nil {
		return err
	}

	// and call image table to get all of image
	images, err := s.store.GetAllImages(destination.ID)
	if err != nil {
		return err
	}

	sendData := &SendSpecificDestinationType{
		ID:               destination.ID,
		Name_Destination: destination.Name_Destination,
		URL_Destination:  destination.URL_Destination,
		List_Image:       images,
	}

	return WriteJSON(w, http.StatusOK, sendData)
}

// handle create new bookmark
func (s *APIServer) handleCreateNewBookmark(w http.ResponseWriter, r *http.Request) error {
	// read data from the body
	book := new(NewBookmarkType)
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		return err
	}

	defer r.Body.Close()

	newBook, err := s.store.CreateNewBookmark(book)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, newBook)
}

// handle get all bookmark name
func (s *APIServer) handleGetBookmarkName(w http.ResponseWriter, r *http.Request) error {
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		return err
	}

	bookmarks, err := s.store.GetAllBookmark(user_id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, bookmarks)
}

// handle save data into bookmark
func (s *APIServer) handleSaveIntoBookmark(w http.ResponseWriter, r *http.Request) error {
	newSave := new(CreateNewUser_SaveType)
	if err := json.NewDecoder(r.Body).Decode(newSave); err != nil {
		return err
	}

	if err := s.store.SaveBookmarkData(newSave); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle create new bookmark and save data
func (s *APIServer) handleCreateAndSaveIntoBookmark(w http.ResponseWriter, r *http.Request) error {
	newBookReq := new(CreateBookmarkAndSaveType)
	if err := json.NewDecoder(r.Body).Decode(newBookReq); err != nil {
		return err
	}

	// create bookmark
	newBookData := &NewBookmarkType{
		User_ID:       newBookReq.User_ID,
		Name_Bookmark: newBookReq.Name_Bookmark,
	}

	newBook, err := s.store.CreateNewBookmark(newBookData)
	if err != nil {
		return err
	}

	// save data
	newSaveData := &CreateNewUser_SaveType{
		Destination_ID: newBookReq.Destination_ID,
		Bookmark_ID:    newBook.ID,
	}

	if err := s.store.SaveBookmarkData(newSaveData); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle update bookmark name
func (s *APIServer) handleBookmarkUpdateName(w http.ResponseWriter, r *http.Request) error {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookmark_id"))
	if err != nil {
		return err
	}

	bookNewName := new(UpdateBookmarkNameType)
	if err := json.NewDecoder(r.Body).Decode(bookNewName); err != nil {
		return err
	}

	// update
	if err := s.store.UpdateBookmarkName(bookID, bookNewName); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle get bookmark data
func (s *APIServer) handleGetBookmarkData(w http.ResponseWriter, r *http.Request) error {
	bookmark_id, err := strconv.Atoi(chi.URLParam(r, "bookmark_id"))
	if err != nil {
		return err
	}

	user_save_data, err := s.store.GetAllDataByBookmark(bookmark_id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user_save_data)
}

// handle delete bookmark name
func (s *APIServer) handleDeleteBookmarkName(w http.ResponseWriter, r *http.Request) error {
	bookmark_id, err := strconv.Atoi(chi.URLParam(r, "bookmark_id"))
	if err != nil {
		return err
	}

	if err := s.store.DeleteBookmark(bookmark_id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// handle delete bookmark destination
func (s *APIServer) handleDeleteBookmarkDestination(w http.ResponseWriter, r *http.Request) error {
	user_dave_id, err := strconv.Atoi(chi.URLParam(r, "destination_book_id"))
	if err != nil {
		return err
	}

	if err := s.store.DeleteBookmarkData(user_dave_id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

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
		c, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				// if cookie is not yet set
				WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "no cookie"})
				return
			}

			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			return
		}

		// get token from jwt
		tokenString := c.Value

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

			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "token invalid"})
			return
		}

		// check is expiration time is end? if end create new one or updated
		if time.Until(claims.ExpiresAt.Time) <= 30*time.Second && time.Until(claims.ExpiresAt.Time) > 0*time.Second {
			newExpirationTime := time.Now().Add(24 * time.Hour)
			claims.ExpiresAt = jwt.NewNumericDate(newExpirationTime)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			newTokenString, err := token.SignedString(jwtKey)

			if err != nil {
				WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   newTokenString,
				Expires: newExpirationTime,
				Domain:  "http://localhost:8080",
				Path:    "/",
			})
		}

		next.ServeHTTP(w, r)
	})
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

// handle send email
func sendEmail(email, token string) error {
	// set email auth
	authEmail := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com")

	// compose email
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" + "Subject: Login Link\r\n" + "\r\n" + "http://127.0.0.1:3000/login/" + token)

	if err := smtp.SendMail("smtp.gmail.com:587", authEmail, os.Getenv("EMAIL"), to, msg); err != nil {
		return err
	}

	return nil
}
