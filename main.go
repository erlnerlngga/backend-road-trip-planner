package main

import (
	"log"
	"os"
)

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
