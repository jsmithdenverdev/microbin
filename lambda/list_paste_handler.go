package lambda

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsmithdenverdev/microbin"
)

type ListPasteHandler struct {
	Store microbin.PasteStore
}

func (l *ListPasteHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pastes, err := microbin.ListPastes(l.Store)(ctx)

	if err != nil {
		return apiGatewayErrorResponse(err, http.StatusInternalServerError), err
	}

	bytes, err := json.Marshal(pastes)

	if err != nil {
		return apiGatewayErrorResponse(err, http.StatusInternalServerError), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(bytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
