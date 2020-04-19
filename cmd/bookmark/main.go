package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"bookmark-api/internal/di"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	bookmarkApi, err := di.CreateBookmarkApi()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(gin.Logger())
	bookmarkApi.RegisterHandlers(r.Group("/"))

	_ = r.Run(":8080")
}
