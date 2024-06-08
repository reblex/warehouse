package main

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	listenAddr string
	storer     Storer
}

func NewServer(listenAddr string, storer Storer) *Server {
	return &Server{
		listenAddr: listenAddr,
		storer:     storer,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("POST /api/products", s.upsert)
	mux.HandleFunc("GET /api/products", s.getAll)

	http.ListenAndServe(":8000", mux)
}

func (s *Server) upsert(w http.ResponseWriter, r *http.Request) {
	var dto ProductsDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storer.Upsert(dto.Products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("products updated"))
}

func (s *Server) getAll(w http.ResponseWriter, r *http.Request) {
	products, err := s.storer.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dto := ProductsDto{
		Products: products,
	}

	json, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
