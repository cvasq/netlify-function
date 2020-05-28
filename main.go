package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func handler2(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 503,
			Body:       "Something went wrong :(",
		}, nil
	}

	cc := lc.ClientContext

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, " + cc.Client.AppTitle,
	}, nil
}

type LocalServer struct{}

func local() {
	server := &LocalServer{}
	fmt.Println("Starting local dev server on :8999")
	http.ListenAndServe(":8999", server)
}

func (l *LocalServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to read request body: %v", err)))
		return
	}

	req := events.APIGatewayProxyRequest{
		Body: string(body),
	}
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	for k, v := range r.Header {
		req.Headers[strings.ToLower(k)] = v[0]
	}
	resp, err := handler2(r.Context(), req)
	if err != nil {
		log.Printf("Error handling request: %v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error handling request: %v", err)))
		return
	}
	for k, v := range resp.Headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(resp.Body))
}

func main2() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		local()
	} else {
		// Make the handler available for Remote Procedure Call by AWS Lambda
		lambda.Start(handler2)
	}
}
