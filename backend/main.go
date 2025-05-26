package main

import (
	"context"
	"fmt"
	"log"
	"urlshortener/db"
	"urlshortener/handlers"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.GetDB().Close(context.Background()); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	fmt.Println("Database initialized successfully")

	handlers.StartServer()
	fmt.Println("Server started successfully")
}
