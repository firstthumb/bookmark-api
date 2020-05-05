package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"bookmark-api/internal/auth"
	"bookmark-api/internal/di"
)

var authorizerApi *auth.Auth

func init() {
	authorizerApi, _ = di.CreateAuth()
}

func Handler(req events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := req.AuthorizationToken
	bearerToken := strings.Split(token, " ")[1]

	claim, err := authorizerApi.VerifyToken(bearerToken)
	if err != nil {
		log.Printf("Unauthorized. Token : %s, Error : %v", token, err)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	methodArn := NewMethodArn(req.MethodArn)
	apiGatewayArn := NewAPIGatewayArn(methodArn.APIGatewayArn)
	principalID := fmt.Sprintf("user|%s", claim.Username)

	resp := NewAuthorizerResponse(principalID, methodArn.AwsAccount)
	resp.Region = methodArn.Region
	resp.APIID = apiGatewayArn.APIID
	resp.Stage = apiGatewayArn.Stage
	// TODO: Be more restrictive
	resp.AllowMethod(All, "/api/*")

	resp.Context = map[string]interface{}{
		"username": claim.Username,
		"method":   claim.Method,
	}

	fmt.Printf("%+v", resp.APIGatewayCustomAuthorizerResponse)

	return resp.APIGatewayCustomAuthorizerResponse, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	lambda.Start(Handler)
}
