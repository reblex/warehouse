package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.ListenAndServe(":8000", mux)
}
