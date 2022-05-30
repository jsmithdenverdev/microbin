package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var (
		requestLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		timingLogger  = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		logger        = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	)

	conf, err := loadConfig()

	if err != nil {
		logger.Fatal(err)
	}

	db, err := connectDB(conf.connection)

	if err != nil {
		logger.Fatal(err)
	}

	s := server{
		config: conf,
		logger: logger,
		router: mux.NewRouter(),
		pasteHandler: pasteHandler{
			logger: logger,
			pasteService: pasteService{
				logger: logger,
				db:     db,
			},
		},
	}

	// configure middleware on the server
	s.router.Use(loggingMiddleware(requestLogger))
	s.router.Use(timingMiddleware(timingLogger))
	s.router.Use(authMiddleware(conf.username, conf.password))

	// configure routes on the server
	s.routes()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.port), s.router); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func connectDB(conn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(conn), &gorm.Config{})

	if err != nil {
		return &gorm.DB{}, err
	}

	if err := db.AutoMigrate(&Paste{}); err != nil {
		return &gorm.DB{}, err
	}

	return db, nil
}
