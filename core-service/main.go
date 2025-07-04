package main

import (
	"log"
)

func main() {
	server, err := InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	server.Run() // listen and serve on 0.0.0.0:8080
}
