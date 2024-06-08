package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printCLIHelp()
		return
	}

	cmd := args[0]
	switch cmd {
	case "ping":
		pingCmd()
	}
}

func printCLIHelp() {
	fmt.Println("Use format: cmd [args]")
	fmt.Println("ping                   - pings all services in docker to check status")
	fmt.Println("articles <filepath>    - uploads articles from given file")
	fmt.Println("products <filepath>    - uploads products from given file")
}

func pingCmd() {
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
