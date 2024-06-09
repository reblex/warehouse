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
	funcArgs := make([]string, 0)
	if len(args) > 1 {
		funcArgs = append(funcArgs, args[1:]...)
	}

	switch cmd {
	case "ping":
		pingCmd()
	case "articles":
		articleCmd(funcArgs[0])
	case "products":
		productCmd(funcArgs[0])
	default:
		printCLIHelp()
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

func articleCmd(filePath string) {
	post("http://localhost:8001/api/articles", filePath)
}

func productCmd(filePath string) {
	post("http://localhost:8002/api/products", filePath)
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

func post(url, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	request, err := http.NewRequest("POST", url, f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer response.Body.Close()

	fmt.Println(response.Status)
}
