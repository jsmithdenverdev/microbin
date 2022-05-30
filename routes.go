package main

func (s *server) routes() {
	pasteRouter := s.router.PathPrefix("/paste").Subrouter()

	// POST /paste
	pasteRouter.
		HandleFunc("", s.pasteHandler.handleCreate()).
		Methods("POST")

	// GET /paste
	pasteRouter.
		HandleFunc("", s.pasteHandler.handleList()).
		Methods("GET")

	// GET /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", s.pasteHandler.handleRead()).
		Methods("GET")

	// DELETE /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", s.pasteHandler.handleDelete()).
		Methods("DELETE")

		// GET /paste/{id}/raw
	pasteRouter.
		HandleFunc("/{id}/raw", s.pasteHandler.handleReadRaw()).
		Methods("GET")
}
