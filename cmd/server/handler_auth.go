package main

import "net/http"

type authHandler struct{}

func (a *authHandler) handleAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	}
}
