package lambda

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

type AuthorizeHandler struct {
	Username string
	Password string
}

func (a *AuthorizeHandler) HandleAPIGateway(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	username, password, ok := basicAuth(event)

	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))

		expectedUsernameHash := sha256.Sum256([]byte(a.Username))
		expectedPasswordHash := sha256.Sum256([]byte(a.Password))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
			return generatePolicy("user", "Allow", event.MethodArn), nil
		}
	}

	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
}
