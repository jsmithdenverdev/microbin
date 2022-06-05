package main

import (
	"context"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	dynamoclient "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jsmithdenverdev/microbin/dynamodb"
	"github.com/jsmithdenverdev/microbin/lambda"
)

func main() {
	conf, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(os.Getenv("AWS_REGION")))

	if err != nil {
		panic(err)
	}

	config, err := loadConfig()

	if err != nil {
		panic(err)
	}

	var (
		client  = dynamoclient.NewFromConfig(conf)
		store   = dynamodb.PasteStore{Client: client, Table: config.tablename}
		handler = lambda.ReadHandler{Store: &store}
	)

	runtime.Start(handler.HandleAPIGateway)
}
