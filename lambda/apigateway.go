package lambda

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

const apiGatewayErrorFormat string = `{"error": "%s"}`

var (
	errNotFound = errors.New("not found")
)

func apiGatewayErrorResponse(e error, status int) events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(apiGatewayErrorFormat, e.Error()),
	}

	// send a 500 status code if an invalid status code is supplied
	if status < 400 || status > 599 {
		response.StatusCode = http.StatusInternalServerError
	} else {
		response.StatusCode = status
	}

	return response
}
