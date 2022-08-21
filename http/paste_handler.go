package http

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

type PasteHandler struct {
	Logger       *log.Logger
	PasteService microbin.PasteService
}

func (p *PasteHandler) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
		case "application/json":
			p.handleCreateText(w, r)
			return
		case "multipart/form-data":
			p.handleCreateFile(w, r)
			return
		default:
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "Content-Type must be multipart/form-data or application/json", http.StatusBadRequest)
			return
		}
	}
}

func (p *PasteHandler) handleCreateText(w http.ResponseWriter, r *http.Request) {
	paste := new(microbin.Paste)

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(paste); err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	paste.Type = microbin.PasteTypeText

	id, err := p.PasteService.Create(*paste)

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(int(id))))
}

func (p *PasteHandler) handleCreateFile(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 1024 * 1024

	paste := new(microbin.Paste)

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

	paste.Type = microbin.PasteTypeFile
	paste.Expiration = microbin.Expiration(r.Header.Get("Expiration"))
	paste.File = fileHeader.Filename
	paste.BinaryContent, err = ioutil.ReadAll(file)

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	id, err := p.PasteService.Create(*paste)

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(int(id))))
}

func (p *PasteHandler) handleRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sID := mux.Vars(r)["id"]

		id, err := strconv.Atoi(sID)

		if err != nil {
			p.Logger.Printf("[ERROR] could not parse id from request: %s\n", err.Error())

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "ID parameter must be an integer.", http.StatusBadRequest)
			return
		}

		paste, err := p.PasteService.Read(id)

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

		enc := json.NewEncoder(w)

		if err := enc.Encode(paste); err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}
	}
}

func (p *PasteHandler) handleReadRaw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sID := mux.Vars(r)["id"]

		id, err := strconv.Atoi(sID)

		if err != nil {
			p.Logger.Printf("[ERROR] could not parse id from request: %s\n", err.Error())

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "ID parameter must be an integer.", http.StatusBadRequest)
			return
		}

		paste, err := p.PasteService.Read(id)

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
		if paste.Type == microbin.PasteTypeText {
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte(paste.Content))
		}

		if paste.Type == microbin.PasteTypeFile {
			ext := filepath.Ext(paste.File)
			mimetype := mime.TypeByExtension(ext)

			w.Header().Set("Content-Type", mimetype)
			w.Write(paste.BinaryContent)
		}
	}
}

func (p *PasteHandler) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sID := mux.Vars(r)["id"]

		id, err := strconv.Atoi(sID)

		if err != nil {
			p.Logger.Printf("[ERROR] could not parse id from request: %s\n", err.Error())

			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, "ID parameter must be an integer.", http.StatusBadRequest)
			return
		}

		err = p.PasteService.Delete(id)

		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(id)))
	}
}

func (p *PasteHandler) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		pastes, err := p.PasteService.List()

		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		enc := json.NewEncoder(w)

		if err := enc.Encode(pastes); err != nil {
			w.Header().Add("Content-Type", "text/plain")
			http.Error(w, errorInternalServer, http.StatusInternalServerError)
		}
	}
}
