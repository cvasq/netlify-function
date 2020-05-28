package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	params := request.QueryStringParameters

	name := "World"
	if params["name"] != "" {
		name = params["name"]
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hi " + name,
	}, nil

}

func main() {
	lambda.Start(handler)
}
