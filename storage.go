package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	SignUp(acc *SignUpType) (*AccountType, error)
	CheckCity(city string) (*CityType, error)
	CreateNewCity(city *CreateNewCityType) (*CityType, error)
	GetDestination(des_id int) (*DestinationType, error)
	CreateDestination(des *CreateNewDestinationType) error
	GetAllDestination(city_id int) ([]*AllDestinationType, error)
}

type MysqlStore struct {
	db *sql.DB
}

func NewPostgresStore() (*MysqlStore, error) {
	// open the connection of db
	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MysqlStore{
		db: db,
	}, nil
}

// crete user tabel
func (s *MysqlStore) CreateTableUser() error {
	createTable := `
		create table if not exists user (
			id int auto_increment,
			name varchar(50) not null,
			email varchar(50) not null unique,
			primary key(id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create table city
func (s *MysqlStore) CreateTableCity() error {
	createTable := `
		create table if not exists city (
			id int auto_increment,
			name_city varchar(50) not null unique,
			lat_city decimal(10,7) not null,
			long_city decimal(10,7) not null,
			primary key(id)
		);
	`
	_, err := s.db.Exec(createTable)

	return err
}

// create table destination
func (s *MysqlStore) CreateTableDestination() error {
	createTable := `
		create table if not exists destination (
			id int auto_increment,
			name_destination varchar(100),
			url_destination varchar(200),
			lat_destination decimal(10,7) not null,
			long_destination decimal(10,7) not null,
			city_id int references city(id),
			primary key(id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create image table
func (s *MysqlStore) CreateTableImage() error {
	createTable := `
		create table if not exists image (
			id int auto_increment,
			url_image varchar(200) not null,
			destination_id int references destination(id),
			primary key(id)
		);
	`
	_, err := s.db.Exec(createTable)

	return err
}

// create bookmark table
func (s *MysqlStore) CreateTableBookmark() error {
	createTable := `
		create table if not exists bookmark (
			id int auto_increment,
			name_bookmark varchar(50) not null,
			user_id int references user(id),
			primary key(id)
		);
	`
	_, err := s.db.Exec(createTable)

	return err
}

// create user_save tabel
func (s *MysqlStore) CreateTableUserSave() error {
	craeteTable := `
		create table if not exists user_save (
			id int auto_increment,
			destination_id int references destination(id),
			bookmark_id int references bookmark(id) on delete cascade,
			primary key(id)
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
	err := s.db.QueryRow(`select * from user where email = $1;`, email).Scan(&acc.ID, &acc.Name, &acc.Email)

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
	

	insertQuery := `insert into user(name, email) values ($1, $2);`

	_, err := s.db.Exec(insertQuery, acc.Name, acc.Email)

	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow(`select * from user where id = (select LAST_INSERT_ID());`).Scan(&account.ID, &account.Name, &account.Email); err != nil {
		return nil, err
	}

	return account, nil
}

// create new city 
func (s *MysqlStore) CreateNewCity(city *CreateNewCityType) (*CityType, error) {
	newCity := new(CityType)

	insertQuery := `insert into city(name, lat_city, long_cty) values ($1, $2, $3);`

	_, err := s.db.Exec(insertQuery, city.Name_City, city.Lat_City, city.Long_City)

	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow(`select * from city where id = (select LAST_INSERT_ID());`).Scan(&newCity.ID, &newCity.Name_City, &newCity.Lat_City, &newCity.Long_City); err != nil {
		return nil, err
	}

	return newCity, err
}

// first check the city its there or not
func (s *MysqlStore) CheckCity(c string) (*CityType, error) {
	city := new(CityType)

	err := s.db.QueryRow("select * from city where name_city = $1;", c).Scan(&city.ID, &city.Name_City, &city.Lat_City, &city.Long_City)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("city: %s not found", c)
	}

	if err != nil {
		return nil , err
	}

	return city, nil
}

// create new destination 
func (s *MysqlStore) CreateNewDestination(des *CreateNewDestinationType) (*DestinationType, error)  {
	newDes := new(DestinationType)

	insertQuery := `insert into destination(name_destination, url_destination, lat_destination, long_destination, city_id) values ($1, $2, $3, $4, $5);`
	
	_, err := s.db.Exec(insertQuery, des.Name_Destination, des.URL_Destination, des.Lat_Destination, des.Long_Destination, des.City_ID)

	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow("select * from destination where id = (select LAST_INSERT_ID());").Scan(&newDes.ID, &newDes.Name_Destination, &newDes.URL_Destination, &newDes.Lat_Destination, &newDes.Long_Destination, &newDes.City_ID); err != nil {
		return nil, err
	}

	return newDes, err
}

// get single image
func (s *MysqlStore) GetSingleImage(des_id int, d *AllDestinationType) (*AllDestinationType, error) {
	err := s.db.QueryRow("select url_image from image where destination_id = $1;", des_id).Scan(&d.URL_Image)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("image id: %d not found", des_id)
	}

	if err != nil {
		return nil, err
	}

	return d, nil
}


// if city is there get all destination data base on city
func (s *MysqlStore) GetAllDestination(city_id int) ([]*AllDestinationType, error)  {
	rows, err := s.db.Query(`select id, name_destination, url_destination, lat_destination, long_destination from destination where city_id = $1;`, city_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allDestination := []*AllDestinationType{}
	for rows.Next() {
		d := new(AllDestinationType)

		if err := rows.Scan(&d.ID, &d.Name_Destination, &d.URL_Destination, &d.Lat_Destination, &d.Long_Destination); err != nil {
			return nil, err
		}

		d, err := s.GetSingleImage(d.ID, d)
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
func (s *MysqlStore) GetDestination(des_id int) (*DestinationType, error) {
	destination := new(DestinationType)

	err := s.db.QueryRow("select * from destination where id = $1;", des_id).Scan(&destination.ID, &destination.Name_Destination, &destination.URL_Destination, &destination.Lat_Destination, &destination.Long_Destination, &destination.City_ID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("destination id: %d not found", des_id)
	}

	if err != nil {
		return nil, err
	}

	return destination, nil
}

// create new images
func (s *MysqlStore) CreateNewImage(img *CreateNewImageType) error {
	insertQuery := `insert into image(url_image, destination_id) values ($1, $2);`

	_, err := s.db.Exec(insertQuery, img.URL_Image, img.Destination_ID)

	if err != nil {
		return err
	}

	return nil
}

// get all Image
func (s *MysqlStore) GetAllImages(des_id int) ([]*ImageType, error)  {
	rows, err := s.db.Query("select * from image where destination_id = $1;", des_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	images := []*ImageType{}
	for rows.Next() {
		i := new(ImageType)

		if err := rows.Scan(&i.ID, &i.URL_Image, &i.Destination_ID); err != nil {
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

	insertQuery := `insert into bookmark(name_bookmark, user_id) values ($1, $2);`

	_, err := s.db.Exec(insertQuery, book.Name_Bookmark, book.User_ID)
	
	if err != nil {
		return nil, err
	}

	if err := s.db.QueryRow("select * from bookmark where id = (select LAST_INSERT_ID());").Scan(&newBook.ID, &newBook.Name_Bookmark, &newBook.User_ID); err != nil {
		return nil , err
	}

	return newBook, err
}


// get all bookmark
func  (s *MysqlStore) GetAllBookmark(user_id int) ([]*BookmarkType, error)  {
	rows, err := s.db.Query("select * from bookmark where user_id = $1;", user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bookmarks := []*BookmarkType{}
	for rows.Next() {
		b := new(BookmarkType)

		if err := rows.Scan(&b.ID, &b.Name_Bookmark, &b.User_ID); err != nil {
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
	insertQuery := `insert into user_save(destination_id, bookmark_id) values($1, $2);`

	_, err := s.db.Exec(insertQuery, newSave.Destination_ID, newSave.Bookmark_ID)

	if err != nil {
		return err
	}

	return nil
}

// get single image
func (s *MysqlStore) GetSingleImageSave_User(des_id int, d *SendDataUser_SaveType) (*SendDataUser_SaveType, error) {
	err := s.db.QueryRow("select url_image from image where destination_id = $1;", des_id).Scan(&d.URL_Image)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("image id: %d not found", des_id)
	}

	if err != nil {
		return nil, err
	}

	return d, nil
}

// get all data from bookmark
func (s *MysqlStore) GetAllDataByBookmark(bookmark_id int) ([]*SendDataUser_SaveType, error) {
	queryStr := "select user_save.id as `user_save_id`, destination.id as `destination_id`, destination.name_destination as `name_destination`, destination.url_destination as `url_destination` from user_save inner join destination on user_save.destination_id = destination.id where user_save.bookmark_id = $1;"

	rows, err := s.db.Query(queryStr, bookmark_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user_save_data := []*SendDataUser_SaveType{}
	for rows.Next() {
		u := new(SendDataUser_SaveType)

		if err := rows.Scan(&u.User_Save_ID, &u.Destination_ID, &u.Name_Destination, &u.URL_Destination); err != nil {
			return nil, err
		}

		u, err := s.GetSingleImageSave_User(u.Destination_ID, u)

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
func (s *MysqlStore) UpdateBookmarkName(id int, name *UpdateBookmarkNameType) error {
	updateQuery := `update bookmark set name_bookmark = $1 where id = $2;`

	_, err := s.db.Exec(updateQuery, name.Name_Bookmark, id)

	if err != nil {
		return err
	}

	return nil
}

// delete bookmark
func (s *MysqlStore) DeleteBookmark(id int) error {
	_, err := s.db.Exec("delete from bookmark where id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

// delete bookmark
func (s *MysqlStore) DeleteBookmarkData(id int) error {
	_, err := s.db.Exec("delete from user_save where id = $1", id)

	if err != nil {
		return err
	}

	return nil
}