package main

import (
	"github.com/golang-jwt/jwt/v4"
)

type SignUpType struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountType struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignInType struct {
	Email string `json:"email"`
}

type ClaimsType struct {
	User_ID int `json:"user_id"`
	jwt.RegisteredClaims
}

// to Get city table
type CityType struct {
	ID        int     `json:"id"`
	Name_City string  `json:"name_city"`
	Lat_City  float64 `json:"lat_city"`
	Long_City float64 `json:"long_city"`
}

// create new city\
type CreateNewCityType struct {
	Name_City string  `json:"name_city"`
	Lat_City  float64 `json:"lat_city"`
	Long_City float64 `json:"long_city"`
}

// to get destination table
type DestinationType struct {
	ID               int     `json:"id"`
	Name_Destination string  `json:"name_destination"`
	URL_Destination  string  `json:"url_destination"`
	Lat_Destination  float64 `json:"lat_destination"`
	Long_Destination float64 `json:"long_destination"`
	City_ID          int     `json:"city_id"`
}

type CreateNewDestinationType struct {
	Name_Destination string  `json:"name_destination"`
	URL_Destination  string  `json:"url_destination"`
	Lat_Destination  float64 `json:"lat_destination"`
	Long_Destination float64 `json:"long_destination"`
	City_ID          int     `json:"city_id"`
}

// send data get All Destination
type AllDestinationType struct {
	ID               int     `json:"id"`
	Name_Destination string  `json:"name_destination"`
	URL_Destination  string  `json:"url_destination"`
	Lat_Destination  float64 `json:"lat_destination"`
	Long_Destination float64 `json:"long_destination"`
	URL_Image        string  `json:"url_image"`
}

type SendAllDestinationType struct {
	Name_City        string                `json:"name_city"`
	Lat_City         float64               `json:"lat_city"`
	Long_City        float64               `json:"long_city"`
	List_Destination []*AllDestinationType `json:"list_destination"`
}

type SendSpecificDestinationType struct {
	ID               int          `json:"id"`
	Name_Destination string       `json:"name_destination"`
	URL_Destination  string       `json:"url_destination"`
	List_Image       []*ImageType `json:"list_image"`
}

type CreateNewImageType struct {
	URL_Image      string `json:"url_image"`
	Destination_ID int    `json:"destination_id"`
}

// to get Image table
type ImageType struct {
	ID             int    `json:"id"`
	URL_Image      string `json:"url_image"`
	Destination_ID int    `json:"destination_id"`
}

type NewBookmarkType struct {
	User_ID       int    `json:"user_id"`
	Name_Bookmark string `json:"name_bookmark"`
}

type CreateBookmarkAndSaveType struct {
	User_ID        int    `json:"user_id"`
	Name_Bookmark  string `json:"name_bookmark"`
	Destination_ID int    `json:"destination_id"`
}

// to get bookmark table
type BookmarkType struct {
	ID            int    `json:"id"`
	Name_Bookmark string `json:"name_bookmark"`
	User_ID       int    `json:"user_id"`
}

// to get user_save tabel
type User_SaveType struct {
	ID               int    `json:"id"`
	Name_Destination string `json:"name_destination"`
	URL_Destination  string `json:"url_destination"`
}

type CreateNewUser_SaveType struct {
	Destination_ID int `json:"destination_id"`
	Bookmark_ID    int `json:"bookmark_id"`
}

type SendDataUser_SaveType struct {
	User_Save_ID     int    `json:"user_save_id"`
	Destination_ID   int    `json:"destination_id"`
	Name_Destination string `json:"name_destination"`
	URL_Destination  string `json:"url_destination"`
	URL_Image        string `json:"url_image"`
}

type UpdateBookmarkNameType struct {
	Name_Bookmark string `json:"name_bookmark"`
}
