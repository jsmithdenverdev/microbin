package main

import (
	"errors"
	"os"
	"strconv"
)

type config struct {
	connection string
	port       int
	username   string
	password   string
}

func loadConfig() (config, error) {
	var (
		portEnv       = os.Getenv("PORT")
		connectionEnv = os.Getenv("CONNECTION_STRING")
		usernameEnv   = os.Getenv("AUTH_USERNAME")
		passwordEnv   = os.Getenv("AUTH_PASSWORD")
	)

	if portEnv == "" {
		return config{}, errors.New("missing required environment variable PORT")
	}

	if connectionEnv == "" {

		return config{}, errors.New("missing required environment variable CONNECTION_STRING")
	}

	if usernameEnv == "" {
		return config{}, errors.New("missing required environment variable AUTH_USERNAME")
	}

	if passwordEnv == "" {

		return config{}, errors.New("missing required environment variable AUTH_PASSWORD")
	}

	port, err := strconv.Atoi(portEnv)

	if err != nil {
		return config{}, err
	}

	return config{
		connection: connectionEnv,
		port:       port,
		username:   usernameEnv,
		password:   passwordEnv,
	}, nil
}
