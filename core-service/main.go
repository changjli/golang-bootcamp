package main

import (
	"log"
)

func main() {
	// db, err := InitializeMigrator()
	// if err != nil {
	// 	log.Fatalf("Could not initialize database for migration: %v", err)
	// }

	// RunMigrations(db.GetInstance())

	server, err := InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	server.Run(":8080") // listen and serve on 0.0.0.0:8080
}
