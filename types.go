package main

import (
	"github.com/golang-jwt/jwt/v4"
)

type SignUpType struct {
	User_Name string `json:"user_name"`
	Email     string `json:"email"`
}

type AccountType struct {
	User_ID   string `json:"user_id"`
	User_Name string `json:"user_name"`
	Email     string `json:"email"`
}

type SignInType struct {
	Email string `json:"email"`
}

type ClaimsType struct {
	User_ID string `json:"user_id"`
	jwt.RegisteredClaims
}

// to Get city table
type CityType struct {
	City_ID   string  `json:"city_id"`
	City_Name string  `json:"city_name"`
	City_Lat  float64 `json:"city_lat"`
	City_Long float64 `json:"city_long"`
}

// create new city\
type CreateNewCityType struct {
	City_Name string  `json:"city_name"`
	City_Lat  float64 `json:"city_lat"`
	City_Long float64 `json:"city_long"`
}

// to get destination table
type DestinationType struct {
	Destination_ID   string  `json:"destination_id"`
	Destination_Name string  `json:"destination_name"`
	Destination_URL  string  `json:"destination_url"`
	Destination_Lat  float64 `json:"destination_lat"`
	Destination_Long float64 `json:"destination_long"`
	City_ID          string  `json:"city_id"`
}

type CreateNewDestinationType struct {
	Destination_Name string  `json:"destination_name"`
	Destination_URL  string  `json:"destination_url"`
	Destination_Lat  float64 `json:"destination_lat"`
	Destination_Long float64 `json:"destination_long"`
	City_ID          string  `json:"city_id"`
}

// send data get All Destination
type AllDestinationType struct {
	Destination_ID   string  `json:"destination_id"`
	Destination_Name string  `json:"destination_name"`
	Destination_URL  string  `json:"destination_url"`
	Destination_Lat  float64 `json:"destination_lat"`
	Destination_Long float64 `json:"destination_long"`
	Image_URL        string  `json:"image_url"`
}

type SendAllDestinationType struct {
	City_Name        string                `json:"city_name"`
	City_Lat         float64               `json:"city_lat"`
	City_Long        float64               `json:"city_long"`
	List_Destination []*AllDestinationType `json:"list_destination"`
}

type SendSpecificDestinationType struct {
	Destination_ID   string       `json:"destination_id"`
	Destination_Name string       `json:"destination_name"`
	Destination_URL  string       `json:"destination_url"`
	List_Image       []*ImageType `json:"list_image"`
}

type CreateNewImageType struct {
	Image_URL      string `json:"image_url"`
	Destination_ID string `json:"destination_id"`
}

// to get Image table
type ImageType struct {
	Image_ID       string `json:"image_id"`
	Image_URL      string `json:"image_url"`
	Destination_ID string `json:"destination_id"`
}

type NewBookmarkType struct {
	User_ID       string `json:"user_id"`
	Bookmark_Name string `json:"bookmark_name"`
}

type CreateBookmarkAndSaveType struct {
	User_ID        string `json:"user_id"`
	Bookmark_Name  string `json:"bookmark_name"`
	Destination_ID string `json:"destination_id"`
}

// to get bookmark table
type BookmarkType struct {
	Bookmark_ID   string `json:"bookmark_id"`
	Bookmark_Name string `json:"bookmark_name"`
	User_ID       string `json:"user_id"`
}

// to get user_save tabel
type User_SaveType struct {
	User_Save_ID     string `json:"user_save_id"`
	Destination_Name string `json:"destination_name"`
	Destination_URL  string `json:"destination_url"`
}

type CreateNewUser_SaveType struct {
	Destination_ID string `json:"destination_id"`
	Bookmark_ID    string `json:"bookmark_id"`
}

type SendDataUser_SaveType struct {
	City_Name        string `json:"city_name"`
	City_ID          string `json:"city_id"`
	User_Save_ID     string `json:"user_save_id"`
	Destination_ID   string `json:"destination_id"`
	Destination_Name string `json:"destination_name"`
	Destination_URL  string `json:"destination_url"`
	Image_URL        string `json:"image_url"`
}

type UpdateBookmarkNameType struct {
	Bookmark_Name string `json:"bookmark_name"`
}
