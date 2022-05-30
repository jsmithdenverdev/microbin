package main

import (
	"log"

	"github.com/gorilla/mux"
)

type server struct {
	config       config
	logger       *log.Logger
	router       *mux.Router
	pasteHandler pasteHandler
}
