package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	mux.HandleFunc(("POST /api/orders"), s.order)

	http.ListenAndServe(":8000", mux)
}

func (s *Server) order(w http.ResponseWriter, r *http.Request) {
	// TODO: Should return the products and articles required? Forgot about that, but the check is valid.
	// TODO If time could return something like this.
	// {
	// 	"products": [
	// 		{
	// 			"Id": "1",
	// 			"Name": "Name",
	// 			"Amount": "3",
	// 			"Articles": [
	// 				{
	// 					"Id": "5",
	// 					"Name": "Name"
	// 					"Amount": "15" // 5 per product = 15 total
	// 				}
	// 			]
	// 		}
	// 	]
	// }
	var dto OrderDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = reserveProducts(dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not complete order: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.Write([]byte("order completed"))
}

func reserveProducts(order OrderDto) error {
	json, err := json.Marshal(order)
	if err != nil {
		return err
	}

	url := "http://product-service:8000/api/products/reserve" // TODO: ENV
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
