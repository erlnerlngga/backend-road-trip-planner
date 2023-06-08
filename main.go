package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	store, err := NewMysqlStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	defer store.db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := NewApiServer("0.0.0.0:"+port, store)
	server.Run()
}
