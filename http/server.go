package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Config       Config
	Logger       *log.Logger
	Router       *mux.Router
	PasteHandler PasteHandler
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(rw, req)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}
