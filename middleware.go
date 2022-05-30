package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func authMiddleware(expectedUsername, expectedPassword string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if ok {
				usernameHash := sha256.Sum256([]byte(username))
				passwordHash := sha256.Sum256([]byte(password))
				expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
				expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

				usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
				passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

				if usernameMatch && passwordMatch {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

func loggingMiddleware(info *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
			info.Printf("%s %s\n", r.Method, r.URL.Path)

			next.ServeHTTP(w, r)
		}))
	}
}

func timingMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start)

			logger.Printf("%s %s took %f seconds\n", r.Method, r.URL.Path, duration.Seconds())
		})
	}
}
