package main

import (
	"log"
)

func main() {
	server, err := InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	server.Run(":8081") // listen and serve on 0.0.0.0:8081
}
