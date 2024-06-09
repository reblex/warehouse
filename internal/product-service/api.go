package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
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
	mux.HandleFunc("POST /api/products/reserve", s.reserve)

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

	products, err = applyAvailability(products)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not apply availability: %s", err.Error()), http.StatusInternalServerError)
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

func (s *Server) reserve(w http.ResponseWriter, r *http.Request) {
	var dto OrderDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	articleReservations := make(map[ArticleId]int)
	for _, order := range dto.Orders {
		product, err := s.storer.GetByName(order.ProductName)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid product name \"%s\"", order.ProductName), http.StatusBadRequest)
			return
		}

		for _, article := range product.Articles {
			articleReservations[article.Id] += article.Amount * order.Amount
		}
	}

	err = reserveArticles(articleReservations)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not reserve products: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.Write([]byte("products reserved"))
}

func reserveArticles(articles map[ArticleId]int) error {
	reservations := make([]ArticleReservation, 0, len(articles))
	for id, count := range articles {
		reservations = append(reservations, ArticleReservation{
			Id:    id,
			Count: count,
		})
	}

	dto := ArticleReservationsDto{
		Reservations: reservations,
	}

	json, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	url := "http://article-service:8000/api/articles/reserve"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(string(body))
	}

	return nil
}

func applyAvailability(products []Product) ([]Product, error) {
	var wg sync.WaitGroup
	updatedProducts := make([]Product, len(products))
	for i, p := range products {
		wg.Add(1)
		go func(idx int, product Product) {
			defer wg.Done()
			availability, _ := getAvailability(articlesToReservation(product.Articles))
			updatedProducts[idx] = Product{
				Name:         product.Name,
				Articles:     product.Articles,
				Price:        product.Price,
				Availability: availability,
			}
		}(i, p)
	}
	wg.Wait()

	return updatedProducts, nil
}

func getAvailability(reservations ArticleReservationsDto) (int, error) {
	jsonData, err := json.Marshal(reservations)
	if err != nil {
		return 0, err
	}

	url := "http://article-service:8000/api/articles/availability"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(string(body))
	}

	var result AvailabilityDto
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to unmarshal json:%s", err.Error())
	}

	return result.Availability, nil
}

func articlesToReservation(articles []Article) ArticleReservationsDto {
	reservations := make([]ArticleReservation, 0, len(articles))
	for _, a := range articles {
		reservations = append(reservations, ArticleReservation{
			Id:    a.Id,
			Count: a.Amount,
		})
	}
	return ArticleReservationsDto{
		Reservations: reservations,
	}
}
