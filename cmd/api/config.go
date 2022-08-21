package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/jsmithdenverdev/microbin/http"
)

func loadConfig() (http.Config, error) {
	var (
		portEnv       = os.Getenv("PORT")
		connectionEnv = os.Getenv("CONNECTION_STRING")
		usernameEnv   = os.Getenv("AUTH_USERNAME")
		passwordEnv   = os.Getenv("AUTH_PASSWORD")
	)

	if portEnv == "" {
		return http.Config{}, errors.New("missing required environment variable PORT")
	}

	if connectionEnv == "" {

		return http.Config{}, errors.New("missing required environment variable CONNECTION_STRING")
	}

	if usernameEnv == "" {
		return http.Config{}, errors.New("missing required environment variable AUTH_USERNAME")
	}

	if passwordEnv == "" {

		return http.Config{}, errors.New("missing required environment variable AUTH_PASSWORD")
	}

	port, err := strconv.Atoi(portEnv)

	if err != nil {
		return http.Config{}, err
	}

	return http.Config{
		Connection: connectionEnv,
		Port:       port,
		Username:   usernameEnv,
		Password:   passwordEnv,
	}, nil
}
