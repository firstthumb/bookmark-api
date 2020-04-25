package main

import (
	"bookmark-api/internal/auth"
	"bookmark-api/internal/di"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

var authorizerApi *auth.Auth

func init() {
	authorizerApi, _ = di.CreateAuth()
}

func generatePolicy(principalId, effect, resource string, username string, method string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	authResponse.Context = map[string]interface{}{
		"username": username,
		"method":   method,
	}

	return authResponse
}

func Handler(req events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := req.AuthorizationToken
	bearerToken := strings.Split(token, " ")[1]

	claim, err := authorizerApi.VerifyToken(bearerToken)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	return generatePolicy("user", "Allow", req.MethodArn, claim.Username, claim.Method), nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	lambda.Start(Handler)
}
