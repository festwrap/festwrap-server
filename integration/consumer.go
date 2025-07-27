package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()

	// Create a Pub/Sub client
	projectID := os.Getenv("FESTWRAP_PUBSUB_PROJECT_ID")
	if projectID == "" {
		log.Fatal("FESTWRAP_PUBSUB_PROJECT_ID environment variable must be set")
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	subscriptionID := os.Getenv("FESTWRAP_PUBSUB_TEST_SUBSCRIPTION")
	if subscriptionID == "" {
		log.Fatal("FESTWRAP_PUBSUB_TEST_SUBSCRIPTION environment variable must be set")
	}

	sub := client.Subscription(subscriptionID)

	// Create a channel to handle signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that will be canceled when we receive a signal
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start receiving messages
	log.Printf("Starting to receive messages from subscription: %s", subscriptionID)
	go func() {
		err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			fmt.Printf("Received message: %s\n", string(msg.Data))
			msg.Ack()
		})
		if err != nil && err != context.Canceled {
			log.Printf("Error receiving messages: %v", err)
		}
	}()

	// Wait for signal to shutdown
	<-signalChan
	log.Println("Shutting down...")
}
