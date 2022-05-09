package main

import "net/http"

func (s *server) routes() {
	pasteRouter := s.router.PathPrefix("/paste").Subrouter()

	// POST /paste
	pasteRouter.
		HandleFunc("", s.handleCreatePaste()).
		Methods("POST")

	// GET /paste
	pasteRouter.
		HandleFunc("", s.handleListPastes()).
		Methods("GET")

	// GET /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", s.handleReadPaste()).
		Methods("GET")

	// DELETE /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", s.handleDeletePaste()).
		Methods("DELETE")

		// GET /paste/{id}/raw
	pasteRouter.
		HandleFunc("/{id}/raw", s.handleReadRawPaste()).
		Methods("GET")
}

func (s *server) middleware() {
	// Logging
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
			s.infoLog.Printf("%s %s\n", r.Method, r.URL.Path)

			next.ServeHTTP(w, r)
		}))
	})
}
