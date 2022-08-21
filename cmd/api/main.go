package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/jsmithdenverdev/microbin"
	"github.com/jsmithdenverdev/microbin/http"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var (
		requestLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		logger        = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	)

	conf, err := loadConfig()

	if err != nil {
		logger.Fatal(err)
	}

	db, err := connectDB(conf.Connection)

	if err != nil {
		logger.Fatal(err)
	}

	s := http.Server{
		Config: conf,
		Logger: logger,
		Router: mux.NewRouter(),
		PasteHandler: http.PasteHandler{
			Logger: logger,
			PasteService: microbin.PasteService{
				DB: db,
			},
		},
	}

	// configure middleware on the server
	s.Router.Use(http.LoggingMiddleware(requestLogger))
	s.Router.Use(http.AuthMiddleware(conf.Username, conf.Password))

	// configure routes on the server
	s.Routes()

	if err := s.ListenAndServe(fmt.Sprintf(":%d", conf.Port)); err != nil {
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
