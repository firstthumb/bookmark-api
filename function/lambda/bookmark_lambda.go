package main

import (
	"bookmark-api/internal/di"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambda

func init() {
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

	bookmarkApi.RegisterHandlers(r.Group("/api/v1", authApi.GetAuthMiddleware().MiddlewareFunc()))

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	lambda.Start(Handler)
}
