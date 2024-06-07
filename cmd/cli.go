package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	pingService("http://localhost:8001/api/ping")
	pingService("http://localhost:8002/api/ping")
	pingService("http://localhost:8003/api/ping")
}

func pingService(url string) {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("%s -> %s\n", url, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s -> %s\n", url, err.Error())
			return
		}
		bodyString := string(bodyBytes)
		fmt.Printf("%s -> %s\n", url, bodyString)
	}
}
