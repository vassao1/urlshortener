package main

import (
	"context"
	"fmt"
	"urlshortener/db"
)

func main() {
	conn, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to the database successfully!")
	fmt.Println(db.GetUrl(conn, "bFGvN4c0"))
}
