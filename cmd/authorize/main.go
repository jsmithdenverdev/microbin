package main

import (
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/jsmithdenverdev/microbin/lambda"
)

func main() {
	config, err := loadConfig()

	if err != nil {
		panic(err)
	}

	handler := lambda.AuthorizeHandler{
		Username: config.username,
		Password: config.password,
	}

	runtime.Start(handler.HandleAPIGateway)
}
