package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func writeJson(w http.ResponseWriter, status int, val any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(val)
}

type APIServer struct {
	listenAddr string
	db         *sql.DB
}

func (s *APIServer) run() {
	http.HandleFunc("GET /albums", s.getAlbumsHandler)
	http.HandleFunc("GET /albums/{id}", s.getAlbumHandler)

	log.Println("API server running at", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, nil))
}

func (s *APIServer) getAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	albums, err := albums(s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, albums)
}

func (s *APIServer) getAlbumHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	album, err := albumByID(s.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, album)
}