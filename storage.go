package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Storage interface {
	CheckEmail(email string) (*AccountType, error)
	SignUp(acc *SignUpType) (*AccountType, error)
	CreateNewCity(city *CreateNewCityType) (*CityType, error)
	CheckCity(c string) (*CityType, error)
	CreateNewDestination(des *CreateNewDestinationType) (*DestinationType, error)
	GetSingleImage(des_id string, d *AllDestinationType) (*AllDestinationType, error)
	GetAllDestination(city_id string) ([]*AllDestinationType, error)
	GetDestination(des_id string) (*DestinationType, error)
	CreateNewImage(img *CreateNewImageType) error
	GetAllImages(des_id string) ([]*ImageType, error)
	CreateNewBookmark(book *NewBookmarkType) (*BookmarkType, error)
	GetAllBookmark(user_id string) ([]*BookmarkType, error)
	SaveBookmarkData(newSave *CreateNewUser_SaveType) error
	GetSingleImageSave_User(des_id string, d *SendDataUser_SaveType) (*SendDataUser_SaveType, error)
	GetAllDataByBookmark(bookmark_id string) ([]*SendDataUser_SaveType, error)
	UpdateBookmarkName(bookmark_id string, name *UpdateBookmarkNameType) error
	DeleteBookmark(bookmark_id string) error
	DeleteBookmarkData(user_save_id string) error
}

type MysqlStore struct {
	db *sql.DB
}

func NewMysqlStore() (*MysqlStore, error) {
	// open the connection of db
	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("database is running...")

	return &MysqlStore{
		db: db,
	}, nil
}

// crete user tabel
func (s *MysqlStore) CreateTableUser() error {
	createTable := `
		create table if not exists user (
			user_id varchar(100),
			user_name varchar(100) not null,
			email varchar(100) not null unique,
			primary key(user_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create table city
func (s *MysqlStore) CreateTableCity() error {
	createTable := `
		create table if not exists city (
			city_id varchar(100),
			city_name varchar(50) not null unique,
			city_lat decimal(10,7) not null,
			city_long decimal(10,7) not null,
			primary key(city_id)
		);
	`
	_, err := s.db.Exec(createTable)

	return err
}

// create table destination
func (s *MysqlStore) CreateTableDestination() error {
	createTable := `
		create table if not exists destination (
			destination_id varchar(100),
			destination_name varchar(100),
			destination_url varchar(200),
			destination_lat decimal(10,7) not null,
			destination_long decimal(10,7) not null,
			city_id varchar(100) references city(city_id),
			primary key(destination_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create image table
func (s *MysqlStore) CreateTableImage() error {
	createTable := `
		create table if not exists image (
			image_id varchar(100),
			image_url varchar(500) not null,
			destination_id varchar(100) references destination(destination_id),
			primary key(image_id)
		);
	`
	_, err := s.db.Exec(createTable)

	return err
}

// create bookmark table
func (s *MysqlStore) CreateTableBookmark() error {
	createTable := `
		create table if not exists bookmark (
			bookmark_id varchar(100),
			bookmark_name varchar(50) not null,
			user_id varchar(100) references user(user_id),
			primary key(bookmark_id)
		);
	`
	_, err := s.db.Exec(createTable)

	return err
}

// create user_save tabel
func (s *MysqlStore) CreateTableUserSave() error {
	craeteTable := `
		create table if not exists user_save (
			user_save_id varchar(100),
			destination_id varchar(100) references destination(destination_id),
			bookmark_id varchar(100) references bookmark(bookmark_id) on delete cascade,
			primary key(user_save_id)
		);
	`

	_, err := s.db.Exec(craeteTable)

	return err
}

func (s *MysqlStore) init() error {

	if err := s.CreateTableUser(); err != nil {
		return err
	}

	if err := s.CreateTableCity(); err != nil {
		return err
	}

	if err := s.CreateTableDestination(); err != nil {
		return err
	}

	if err := s.CreateTableImage(); err != nil {
		return err
	}

	if err := s.CreateTableBookmark(); err != nil {
		return err
	}

	if err := s.CreateTableUserSave(); err != nil {
		return err
	}

	return nil
}

// check email
func (s *MysqlStore) CheckEmail(email string) (*AccountType, error) {
	acc := new(AccountType)
	err := s.db.QueryRow(`select * from user where email = ?;`, email).Scan(&acc.User_ID, &acc.User_Name, &acc.Email)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account %s not found", email)
	}

	if err != nil {
		return nil, err
	}

	return acc, nil
}

// Sign Up
func (s *MysqlStore) SignUp(acc *SignUpType) (*AccountType, error) {
	account := new(AccountType)

	id := uuid.New().String()

	insertQuery := `insert into user(user_id, user_name, email) values (?, ?, ?);`

	_, err := s.db.Exec(insertQuery, id, acc.User_Name, acc.Email)

	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow(`select * from user where user_id = ?;`, id).Scan(&account.User_ID, &account.User_Name, &account.Email); err != nil {
		return nil, err
	}

	return account, nil
}

// create new city
func (s *MysqlStore) CreateNewCity(city *CreateNewCityType) (*CityType, error) {
	newCity := new(CityType)

	id := uuid.New().String()

	insertQuery := `insert into city(city_id, city_name, city_lat, city_long) values (?, ?, ?, ?);`

	_, err := s.db.Exec(insertQuery, id, city.City_Name, city.City_Lat, city.City_Long)

	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow(`select * from city where city_id = ?;`, id).Scan(&newCity.City_ID, &newCity.City_Name, &newCity.City_Lat, &newCity.City_Long); err != nil {
		return nil, err
	}

	return newCity, err
}

// first check the city its there or not
func (s *MysqlStore) CheckCity(c string) (*CityType, error) {
	city := new(CityType)

	err := s.db.QueryRow("select * from city where city_name = ?;", c).Scan(&city.City_ID, &city.City_Name, &city.City_Lat, &city.City_Long)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("city: %s not found", c)
	}

	if err != nil {
		return nil, err
	}

	return city, nil
}

// create new destination
func (s *MysqlStore) CreateNewDestination(des *CreateNewDestinationType) (*DestinationType, error) {
	newDes := new(DestinationType)

	id := uuid.New().String()

	insertQuery := `insert into destination(destination_id, destination_name, destination_url, destination_lat, destination_long, city_id) values (?, ?, ?, ?, ?, ?);`

	_, err := s.db.Exec(insertQuery, id, des.Destination_Name, des.Destination_URL, des.Destination_Lat, des.Destination_Long, des.City_ID)

	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow("select * from destination where destination_id = ?;", id).Scan(&newDes.Destination_ID, &newDes.Destination_Name, &newDes.Destination_URL, &newDes.Destination_Lat, &newDes.Destination_Long, &newDes.City_ID); err != nil {
		return nil, err
	}

	return newDes, err
}

// get single image
func (s *MysqlStore) GetSingleImage(des_id string, d *AllDestinationType) (*AllDestinationType, error) {
	err := s.db.QueryRow("select image_url from image where destination_id = ?;", des_id).Scan(&d.Image_URL)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("image id: %s not found", des_id)
	}

	if err != nil {
		return nil, err
	}

	return d, nil
}

// if city is there get all destination data base on city
func (s *MysqlStore) GetAllDestination(city_id string) ([]*AllDestinationType, error) {
	rows, err := s.db.Query(`select destination_id, destination_name, destination_url, destination_lat, destination_long from destination where city_id = ?;`, city_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allDestination := []*AllDestinationType{}
	for rows.Next() {
		d := new(AllDestinationType)

		if err := rows.Scan(&d.Destination_ID, &d.Destination_Name, &d.Destination_URL, &d.Destination_Lat, &d.Destination_Long); err != nil {
			return nil, err
		}

		d, err := s.GetSingleImage(d.Destination_ID, d)
		if err != nil {
			return nil, err
		}

		allDestination = append(allDestination, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allDestination, err
}

// get single destination
func (s *MysqlStore) GetDestination(des_id string) (*DestinationType, error) {
	destination := new(DestinationType)

	err := s.db.QueryRow("select * from destination where destination_id = ?;", des_id).Scan(&destination.Destination_ID, &destination.Destination_Name, &destination.Destination_URL, &destination.Destination_Lat, &destination.Destination_Long, &destination.City_ID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("destination id: %s not found", des_id)
	}

	if err != nil {
		return nil, err
	}

	return destination, nil
}

// create new images
func (s *MysqlStore) CreateNewImage(img *CreateNewImageType) error {
	id := uuid.New().String()

	insertQuery := `insert into image(image_id, image_url, destination_id) values (?, ?, ?);`

	_, err := s.db.Exec(insertQuery, id, img.Image_URL, img.Destination_ID)

	if err != nil {
		return err
	}

	return nil
}

// get all Image
func (s *MysqlStore) GetAllImages(des_id string) ([]*ImageType, error) {
	rows, err := s.db.Query("select * from image where destination_id = ?;", des_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	images := []*ImageType{}
	for rows.Next() {
		i := new(ImageType)

		if err := rows.Scan(&i.Image_ID, &i.Image_URL, &i.Destination_ID); err != nil {
			return nil, err
		}

		images = append(images, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

// create new bookmark
func (s *MysqlStore) CreateNewBookmark(book *NewBookmarkType) (*BookmarkType, error) {
	newBook := new(BookmarkType)
	id := uuid.New().String()

	insertQuery := `insert into bookmark(bookmark_id, bookmark_name, user_id) values (?, ?, ?);`

	_, err := s.db.Exec(insertQuery, id, book.Bookmark_Name, book.User_ID)

	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow("select * from bookmark where bookmark_id = ?;", id).Scan(&newBook.Bookmark_ID, &newBook.Bookmark_Name, &newBook.User_ID); err != nil {
		return nil, err
	}

	return newBook, err
}

// get all bookmark
func (s *MysqlStore) GetAllBookmark(user_id string) ([]*BookmarkType, error) {
	rows, err := s.db.Query("select * from bookmark where user_id = ?;", user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bookmarks := []*BookmarkType{}
	for rows.Next() {
		b := new(BookmarkType)

		if err := rows.Scan(&b.Bookmark_ID, &b.Bookmark_Name, &b.User_ID); err != nil {
			return nil, err
		}

		bookmarks = append(bookmarks, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookmarks, nil
}

// save bookmark data
func (s *MysqlStore) SaveBookmarkData(newSave *CreateNewUser_SaveType) error {
	id := uuid.New().String()
	insertQuery := `insert into user_save(user_save_id, destination_id, bookmark_id) values(?, ?, ?);`

	_, err := s.db.Exec(insertQuery, id, newSave.Destination_ID, newSave.Bookmark_ID)

	if err != nil {
		return err
	}

	return nil
}

// get single image
func (s *MysqlStore) GetSingleImageSave_User(des_id string, d *SendDataUser_SaveType) (*SendDataUser_SaveType, error) {
	err := s.db.QueryRow("select image_url from image where destination_id = ?;", des_id).Scan(&d.Image_URL)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("image id: %s not found", des_id)
	}

	if err != nil {
		return nil, err
	}

	return d, nil
}

// get city name
func (s *MysqlStore) getCityName(city_id string, d *SendDataUser_SaveType) (*SendDataUser_SaveType, error) {

	err := s.db.QueryRow("select city_name from city where city_id = ?;", city_id).Scan(&d.City_Name)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("city: %s not found", city_id)
	}

	if err != nil {
		return nil, err
	}

	return d, nil
}

// get all data from bookmark
func (s *MysqlStore) GetAllDataByBookmark(bookmark_id string) ([]*SendDataUser_SaveType, error) {
	queryStr := "select user_save.user_save_id as `user_save_id`, destination.destination_id as `destination_id`, destination.destination_name as `destination_name`, destination.destination_url as `destination_url`, destination.city_id as `city_id` from user_save inner join destination on user_save.destination_id = destination.destination_id where user_save.bookmark_id = ?;"

	rows, err := s.db.Query(queryStr, bookmark_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user_save_data := []*SendDataUser_SaveType{}
	for rows.Next() {
		u := new(SendDataUser_SaveType)

		if err := rows.Scan(&u.User_Save_ID, &u.Destination_ID, &u.Destination_Name, &u.Destination_URL, &u.City_ID); err != nil {
			return nil, err
		}

		u, err := s.GetSingleImageSave_User(u.Destination_ID, u)

		if err != nil {
			return nil, err
		}

		u, err = s.getCityName(u.City_ID, u)

		if err != nil {
			return nil, err
		}

		user_save_data = append(user_save_data, u)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return user_save_data, nil
}

// update bookmark name
func (s *MysqlStore) UpdateBookmarkName(bookmark_id string, name *UpdateBookmarkNameType) error {
	updateQuery := `update bookmark set bookmark_name = ? where bookmark_id = ?;`

	_, err := s.db.Exec(updateQuery, name.Bookmark_Name, bookmark_id)

	if err != nil {
		return err
	}

	return nil
}

// delete bookmark
func (s *MysqlStore) DeleteBookmark(bookmark_id string) error {
	_, err := s.db.Exec("delete from user_save where bookmark_id = ?;", bookmark_id)

	if err != nil {
		return err
	}

	_, err = s.db.Exec("delete from bookmark where bookmark_id = ?;", bookmark_id)

	if err != nil {
		return err
	}

	return nil
}

// delete bookmark
func (s *MysqlStore) DeleteBookmarkData(user_save_id string) error {
	_, err := s.db.Exec("delete from user_save where user_save_id = ?", user_save_id)

	if err != nil {
		return err
	}

	return nil
}
