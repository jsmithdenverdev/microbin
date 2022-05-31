package main

import (
	"fmt"
	"os"
	"strings"
)

type config struct {
	username string
	password string
}

func loadConfig() (config, error) {
	missing := []string{}

	var (
		username = os.Getenv("AUTH_USERNAME")
		password = os.Getenv("AUTH_PASSWORD")
	)

	if username == "" {
		missing = append(missing, "AUTH_USERNAME")
	}

	if password == "" {
		missing = append(missing, "AUTH_PASSWORD")
	}

	if len(missing) > 0 {
		if len(missing) > 1 {
			return config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
		} else {
			return config{}, fmt.Errorf("missing required environment variable: %s", missing[0])
		}
	}

	return config{
		username: username,
		password: password,
	}, nil
}
