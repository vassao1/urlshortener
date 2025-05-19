package db

import (
	"context"
	"fmt"
	"os"
	"urlshortener/shortener"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return conn, nil
}

func AddUrl(conn *pgx.Conn, longUrl string) (string, error) {
	shortUrl := shortener.GenerateShortURL()

	_, err := conn.Exec(context.Background(), "INSERT INTO urls (original, shortened) VALUES ($1, $2)", longUrl, shortUrl)
	if err != nil {
		return "", fmt.Errorf("error inserting URL into database: %v", err)
	}

	return shortUrl, nil
}

func GetUrl(conn *pgx.Conn, shortUrl string) (string, error) {
	var longUrl string
	err := conn.QueryRow(context.Background(), "SELECT original FROM urls WHERE shortened = $1", shortUrl).Scan(&longUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("short URL not found")
		}
		return "", fmt.Errorf("error retrieving URL from database: %v", err)
	}

	return longUrl, nil
}
