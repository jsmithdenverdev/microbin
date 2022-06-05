package main

import (
	"fmt"
	"os"
	"strings"
)

type config struct {
	tablename string
}

func loadConfig() (config, error) {
	missing := []string{}

	var (
		tablename = os.Getenv("TABLE_NAME")
	)

	if tablename == "" {
		missing = append(missing, "TABLE_NAME")
	}

	if len(missing) > 0 {
		if len(missing) > 1 {
			return config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
		} else {
			return config{}, fmt.Errorf("missing required environment variable: %s", missing[0])
		}
	}

	return config{
		tablename: tablename,
	}, nil
}
