package main

func (a *application) routes() {
	pasteRouter := a.router.PathPrefix("/paste").Subrouter()

	// POST /paste
	pasteRouter.
		HandleFunc("", a.pasteHandler.handleCreate()).
		Methods("POST")

	// GET /paste
	pasteRouter.
		HandleFunc("", a.pasteHandler.handleList()).
		Methods("GET")

	// GET /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", a.pasteHandler.handleRead()).
		Methods("GET")

	// DELETE /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", a.pasteHandler.handleDelete()).
		Methods("DELETE")

		// GET /paste/{id}/raw
	pasteRouter.
		HandleFunc("/{id}/raw", a.pasteHandler.handleReadRaw()).
		Methods("GET")
}
