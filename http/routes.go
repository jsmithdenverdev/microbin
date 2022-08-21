package http

func (s *Server) Routes() {
	pasteRouter := s.Router.PathPrefix("/paste").Subrouter()

	// POST /paste
	pasteRouter.
		HandleFunc("", s.PasteHandler.handleCreate()).
		Methods("POST")

	// GET /paste
	pasteRouter.
		HandleFunc("", s.PasteHandler.handleList()).
		Methods("GET")

	// GET /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", s.PasteHandler.handleRead()).
		Methods("GET")

	// DELETE /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", s.PasteHandler.handleDelete()).
		Methods("DELETE")

		// GET /paste/{id}/raw
	pasteRouter.
		HandleFunc("/{id}/raw", s.PasteHandler.handleReadRaw()).
		Methods("GET")
}
