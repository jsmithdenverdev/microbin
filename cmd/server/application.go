package main

import (
	"log"

	"github.com/gorilla/mux"
)

type application struct {
	config       config
	infoLog      *log.Logger
	errorLog     *log.Logger
	router       *mux.Router
	authHandler  authHandler
	pasteHandler pasteHandler
}
