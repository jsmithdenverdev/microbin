package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsmithdenverdev/microbin"
)

type CreateHandler struct {
	Store microbin.PasteStore
}

func (c *CreateHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

type ReadHandler struct {
	Store microbin.PasteStore
}

func (r *ReadHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

type DeleteHandler struct {
	Store microbin.PasteStore
}

func (d *DeleteHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

type ListHandler struct {
	Store microbin.PasteStore
}

func (l *ListHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
