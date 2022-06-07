package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsmithdenverdev/microbin"
)

type CreatePasteHandler struct {
	Store microbin.PasteStore
}

func (c *CreatePasteHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	paste := microbin.Paste{}

	dec := json.NewDecoder(bytes.NewBuffer([]byte(event.Body)))
	err := dec.Decode(&paste)

	if err != nil {
		return apiGatewayErrorResponse(err, http.StatusInternalServerError), err
	}

	err = microbin.CreatePaste(c.Store)(ctx, paste)

	if err != nil {
		if errors.Is(err, microbin.ErrInvalidExpiration) {
			return apiGatewayErrorResponse(err, http.StatusBadRequest), err
		}
		return apiGatewayErrorResponse(err, http.StatusInternalServerError), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
