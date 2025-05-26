package shortener

import (
	"fmt"
	"math/rand"
)

func GenerateShortURL() string {
	fmt.Print("generate short url chamado")
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 8
	shortURL := make([]byte, length)
	for i := 0; i < length; i++ {
		index := rand.Intn(len(charset))
		shortURL[i] = charset[index]
	}
	fmt.Print("Generated short URL: ")
	fmt.Println(string(shortURL))
	return string(shortURL)
}
