package main

func (a *application) routes() {
	pasteRouter := a.router.PathPrefix("/paste").Subrouter()

	// POST /paste
	pasteRouter.
		HandleFunc("", a.handleCreatePaste()).
		Methods("POST")

	// GET /paste
	pasteRouter.
		HandleFunc("", a.handleListPastes()).
		Methods("GET")

	// GET /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", a.handleReadPaste()).
		Methods("GET")

	// DELETE /paste/{id}
	pasteRouter.
		HandleFunc("/{id}", a.handleDeletePaste()).
		Methods("DELETE")

		// GET /paste/{id}/raw
	pasteRouter.
		HandleFunc("/{id}/raw", a.handleReadRawPaste()).
		Methods("GET")
}
