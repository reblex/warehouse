package main

func main() {
	server := NewServer(":8000", NewMemoryStorer())
	server.Start()
}
