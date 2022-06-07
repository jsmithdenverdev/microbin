package lambda

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsmithdenverdev/microbin"
)

type DeletePasteHandler struct {
	Store microbin.PasteStore
}

func (d *DeletePasteHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idStr, ok := event.PathParameters["id"]

	if !ok {
		return apiGatewayErrorResponse(errNotFound, http.StatusNotFound), errNotFound
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return apiGatewayErrorResponse(err, http.StatusNotFound), err
	}

	err = microbin.DeletePaste(d.Store)(ctx, id)

	if err != nil {
		return apiGatewayErrorResponse(err, http.StatusInternalServerError), err
	}

	if err != nil {
		return apiGatewayErrorResponse(err, http.StatusInternalServerError), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
