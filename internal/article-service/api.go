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
	mux.HandleFunc("POST /api/articles", s.storeArticles)
	mux.HandleFunc("GET /api/articles", s.getAllArticles)
	mux.HandleFunc("POST /api/articles/reserve", s.reserveArticles)
	mux.HandleFunc("POST /api/articles/availability", s.calculateAvailability)

	http.ListenAndServe(":8000", mux)
}

func (s *Server) storeArticles(w http.ResponseWriter, r *http.Request) {
	var dto ArticlesDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storer.Store(dto.Articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("articles stored successfully"))
}

func (s *Server) getAllArticles(w http.ResponseWriter, _ *http.Request) {
	articles, err := s.storer.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dto := ArticlesDto{
		Articles: articles,
	}

	json, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (s *Server) reserveArticles(w http.ResponseWriter, r *http.Request) {
	var dto ReservationsDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storer.Reserve(dto.Reservations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("articles reserved successfully"))
}

func (s *Server) calculateAvailability(w http.ResponseWriter, r *http.Request) {
	var dto ReservationsDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData := AvailabilityDto{
		Availability: s.storer.CalculateAvailability(dto.Reservations),
	}

	json, err := json.Marshal(resData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
