package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jsmithdenverdev/microbin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)

	conf, err := loadConfig()

	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := connectDB(conf.connection)

	if err != nil {
		errorLog.Fatal(err)
	}

	a := application{
		config:   conf,
		infoLog:  infoLog,
		errorLog: errorLog,
		router:   mux.NewRouter(),
		pasteHandler: pasteHandler{
			infoLog:      infoLog,
			errorLog:     errorLog,
			pasteService: microbin.NewPasteService(db, infoLog, errorLog),
		},
	}

	// configure middleware on the server
	a.middleware()

	// configure routes on the server
	a.routes()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.port), a.router); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func connectDB(conn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(conn), &gorm.Config{})

	if err != nil {
		return &gorm.DB{}, err
	}

	if err := db.AutoMigrate(&microbin.Paste{}); err != nil {
		return &gorm.DB{}, err
	}

	return db, nil
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
