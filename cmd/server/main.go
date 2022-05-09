package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jsmithdenverdev/microbin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	conn := flag.String("conn", "microbin.db", "Sqlite3 database connection string")
	port := flag.Int("port", 8080, "Server port to listen for incoming connections")

	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)

	db, err := connectDB(*conn)

	if err != nil {
		log.Fatal(err)
	}

	s := server{
		infoLog:  infoLog,
		errorLog: errorLog,
		router:   mux.NewRouter(),
		pasteService: microbin.NewPasteService(
			db,
			infoLog,
			errorLog,
		),
	}

	// configure middleware on the server
	s.middleware()

	// configure routes on the server
	s.routes()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), s.router); err != nil && err != http.ErrServerClosed {
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
