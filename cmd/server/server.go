package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jsmithdenverdev/microbin"
	"gorm.io/gorm"
)

type server struct {
	infoLog      *log.Logger
	errorLog     *log.Logger
	router       *mux.Router
	pasteService *microbin.PasteService
}

func (s *server) handleCreatePaste() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
		case "application/json":
			s.handleCreatePasteRaw(w, r)
			return
		case "multipart/form-data":
			s.handleCreatePasteFile(w, r)
			return
		default:
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "Content-Type must be multipart/form-data or application/json", http.StatusBadRequest)
			return
		}
	}
}

func (s *server) handleCreatePasteRaw(w http.ResponseWriter, r *http.Request) {
	paste := microbin.Paste{}
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&paste); err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	id, err := s.pasteService.Create(paste)

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(int(id))))
}

func (s *server) handleCreatePasteFile(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("paste")

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	defer file.Close()

	paste := microbin.Paste{}

	paste.Expiration = r.Header.Get("Expiration")
	paste.File = fileHeader.Filename
	paste.BinaryContent, err = ioutil.ReadAll(file)

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	id, err := s.pasteService.Create(paste)

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(int(id))))
}

func (s *server) handleReadPaste() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sID := mux.Vars(r)["id"]

		id, err := strconv.Atoi(sID)

		if err != nil {
			s.errorLog.Printf("could not parse id from request: %s\n", err.Error())

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "ID parameter must be an integer.", http.StatusBadRequest)
			return
		}

		paste, err := s.pasteService.Read(id)

		if err != nil {
			// FIXME: Don't leak ORM implementation details to the controller (gorm.ErrRecordNotFound)
			if errors.As(err, &microbin.ErrorPasteExpired{}) || errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}

		// TECHDEBT: is there a better convention than this?
		if paste.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)

		if err := enc.Encode(paste); err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleReadRawPaste() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sID := mux.Vars(r)["id"]

		id, err := strconv.Atoi(sID)

		if err != nil {
			s.errorLog.Printf("could not parse id from request: %s\n", err.Error())

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "ID parameter must be an integer.", http.StatusBadRequest)
			return
		}

		paste, err := s.pasteService.Read(id)

		if err != nil {
			// FIXME: Don't leak ORM implementation details to the controller (gorm.ErrRecordNotFound)
			if errors.As(err, &microbin.ErrorPasteExpired{}) || errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}

		// TECHDEBT: is there a better convention than this?
		if paste.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// paste is just text content
		if paste.File == "" {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(paste.Content))
		}

		// TECHDEBT: Create a ContentType on the Paste to avoid doing these checks
		if len(paste.BinaryContent) > 0 {
			ext := filepath.Ext(paste.File)
			mimetype := mime.TypeByExtension(ext)

			w.Header().Set("Content-Type", mimetype)
			w.WriteHeader(http.StatusOK)
			w.Write(paste.BinaryContent)
		}
	}
}

func (s *server) handleDeletePaste() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sID := mux.Vars(r)["id"]

		id, err := strconv.Atoi(sID)

		if err != nil {
			s.errorLog.Printf("could not parse id from request: %s\n", err.Error())

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "ID parameter must be an integer.", http.StatusBadRequest)
			return
		}

		err = s.pasteService.Delete(id)

		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(id)))
	}
}

func (s *server) handleListPastes() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		pastes, err := s.pasteService.List()

		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)

		if err := enc.Encode(pastes); err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
		}
	}
}
