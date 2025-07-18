package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	quit := make(chan os.Signal, 1)

	// Notify the channel on SIGINT (Ctrl+C) or SIGTERM.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	<-quit

	// When a signal is received, the program continues from here.
	log.Println("Shutting down application...")
}
