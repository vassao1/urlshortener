package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"urlshortener/shortener"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var (
	conn     *pgx.Conn
	connOnce sync.Once
)

func InitDB() error {
	var err error
	connOnce.Do(func() {
		err = godotenv.Load()
		if err != nil {
			err = fmt.Errorf("error loading .env file: %v", err)
			return
		}
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			err = fmt.Errorf("DATABASE_URL environment variable not set")
			return
		}
		conn, err = pgx.Connect(context.Background(), dbURL)
		if err != nil {
			err = fmt.Errorf("unable to connect to database: %v", err)
		}
	})
	return err
}

func GetDB() *pgx.Conn {
	return conn
}

func AddUrl(longUrl string) (string, error) {
	fmt.Print("db addurl chamado")
	shortUrl := shortener.GenerateShortURL()
	if shortUrl == "" {
		return "", fmt.Errorf("error generating short URL")
	}
	for !isUnique(shortUrl) {
		shortUrl = shortener.GenerateShortURL()
	}

	_, err := conn.Exec(context.Background(), "INSERT INTO urls (original, shortened) VALUES ($1, $2)", longUrl, shortUrl)
	if err != nil {
		return "", fmt.Errorf("error inserting URL into database: %v", err)
	}

	return shortUrl, nil
}

func GetUrl(shortUrl string) (string, error) {
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

func isUnique(shorturl string) bool {
	url, _ := GetUrl(shorturl)
	if shorturl == "" {
		return false
	}
	if url == "" {
		return true
	} else {
		return false
	}
}
