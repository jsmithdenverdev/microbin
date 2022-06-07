package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsmithdenverdev/microbin"
)

type ReadPasteHandler struct {
	Store microbin.PasteStore
}

func (r *ReadPasteHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idStr, ok := event.PathParameters["id"]

	if !ok {
		return apiGatewayErrorResponse(errNotFound, http.StatusNotFound), errNotFound
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apiGatewayErrorResponse(err, http.StatusNotFound), err
	}

	paste, err := microbin.ReadPaste(r.Store)(ctx, id)

	if err != nil {
		if errors.As(err, microbin.ErrExpiredPaste) {
			return apiGatewayErrorResponse(errNotFound, http.StatusNotFound), err
		}
		return apiGatewayErrorResponse(err, http.StatusInternalServerError), err
	}

	bytes, err := json.Marshal(paste)

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
