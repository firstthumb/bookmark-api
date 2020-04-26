package main

import (
	"log"

	"github.com/gin-contrib/cors"
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

	authApi, err := di.CreateAuthApi()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(cors.Default())
	r.Use(authApi.Session())

	authApi.RegisterSigninHandlers(r.Group("/"))
	authApi.RegisterAuthHandlers(r.Group("/"))

	bookmarkApi.RegisterHandlers(r.Group("/", authApi.GetAuthMiddleware().MiddlewareFunc()))

	_ = r.Run(":8080")
}
