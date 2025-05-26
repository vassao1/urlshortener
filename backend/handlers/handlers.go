package handlers

import (
	"net/http"
	"urlshortener/db"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/generate", func(c *gin.Context) {
		longURL := c.Query("longurl")

		if longURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'long_url' in query"})
			return
		}

		shortURL, err := db.AddUrl(longURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
	})

	r.POST("/shortenlink", func(c *gin.Context) {
		var requestBody struct {
			LongURL string `json:"long_url"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil || requestBody.LongURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'long_url' in request body"})
			return
		}

		if len(requestBody.LongURL) >= 7 && requestBody.LongURL[:7] == "http://" {
		} else if len(requestBody.LongURL) < 8 || requestBody.LongURL[:8] != "https://" {
			requestBody.LongURL = "https://" + requestBody.LongURL
		}

		shortURL, err := db.AddUrl(requestBody.LongURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
	})

	r.GET("/load/:short_url", func(c *gin.Context) {
		shortURL := c.Param("short_url")
		if shortURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'short_url' in path"})
			return
		}
		longURL, err := db.GetUrl(shortURL)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		c.Redirect(http.StatusFound, longURL)
	})

	return r
}

func StartServer() {
	r := InitRoutes()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
