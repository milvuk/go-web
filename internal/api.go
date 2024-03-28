package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	ListenAddr string
	DB         *sql.DB
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *APIServer) Run() {
	http.HandleFunc("GET /albums/", s.getAlbumsHandler)
	http.HandleFunc("GET /albums/{id}", s.getAlbumHandler)
	http.HandleFunc("POST /albums", s.withJWTAuth(s.postAlbumHandler))
	http.HandleFunc("DELETE /albums/{id}", s.withJWTAuth(s.deleteAlbumHandler))
	http.HandleFunc("PUT /albums/{id}", s.withJWTAuth(s.updateAlbumHandler))

	http.HandleFunc("POST /login", s.loginHandler)

	// pass-through mockapi client
	http.HandleFunc("GET /mockapi/products", s.getProductsHandler)
	http.HandleFunc("GET /mockapi/products/{id}", s.getProductHandler)

	log.Println("API server running at", s.ListenAddr)
	log.Fatal(http.ListenAndServe(s.ListenAddr, nil))
}

func (s *APIServer) getAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	// todo: pagination, filtering, ordering
	albums, err := albums(s.DB)
	if err != nil {
		log.Println(err)
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
	album, err := albumByID(s.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, album)
}

func (s *APIServer) postAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var alb Album

	err := json.NewDecoder(r.Body).Decode(&alb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertId, err := addAlbum(s.DB, alb)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/albums/%d", insertId))

	var createdAlb Album
	createdAlb, err = albumByID(s.DB, insertId)
	if err != nil {
		writeJson(w, http.StatusCreated, nil)
		return
	}

	writeJson(w, http.StatusCreated, createdAlb)
}

func (s *APIServer) deleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = deleteAlbum(s.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusNoContent, "")
}

func (s *APIServer) updateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var alb Album

	err = json.NewDecoder(r.Body).Decode(&alb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = updateAlbum(s.DB, id, alb)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var updatedAlb Album
	updatedAlb, err = albumByID(s.DB, id)
	if err != nil {
		writeJson(w, http.StatusOK, nil)
		return
	}

	writeJson(w, http.StatusOK, updatedAlb)
}

func (s *APIServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	adminUser := ViperEnvVariable("ADMIN_USER")
	adminPass := ViperEnvVariable("ADMIN_PASS")

	if u.Username == adminUser && u.Password == adminPass {
		tokenString, err := createToken(u.Username)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		type loginResponse struct {
			Token string `json:"token"`
		}

		writeJson(w, http.StatusOK, &loginResponse{Token: tokenString})
		return
	}
	http.Error(w, "", http.StatusUnauthorized)
}

func (s *APIServer) getProductsHandler(w http.ResponseWriter, r *http.Request) {
	// pass-through mockapi data
	products, err := retrieveProducts()
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, products)
}

func (s *APIServer) getProductHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	product, err := retrieveProduct(id)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, product)
}

func (s *APIServer) withJWTAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Permission denied", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]

		err := verifyToken(tokenString)
		if err != nil {
			log.Println(err)
			http.Error(w, "Permission denied", http.StatusUnauthorized)
			return
		}

		f(w, r)
	}
}