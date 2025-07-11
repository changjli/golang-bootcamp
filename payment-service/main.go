package main

import (
	"context"
	"log"
)

func main() {
	// 1. Initialize the application's components using the injector.
	// The injector now returns a single App struct containing all top-level components.
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	// 2. Create a context that we can cancel on shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Start the RabbitMQ consumer in a separate goroutine.
	// This will listen for new payment requests.
	go app.Consumer.StartConsumer(ctx)
}
