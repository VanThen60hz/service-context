package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	sctx "github.com/VanThen60hz/service-context"
	"github.com/VanThen60hz/service-context/component/natsc"
	"github.com/VanThen60hz/service-context/component/pubsub"
)

func main() {
	// Create service context
	ctx := sctx.NewServiceContext()

	// Initialize NATS component
	natsComp := natsc.NewNatsComponent("nats")
	natsComp.InitFlags()

	// Set NATS server address via flag
	flag.Set("nats-server", "nats://localhost:4222")
	flag.Parse()

	// Activate NATS component
	if err := natsComp.Activate(ctx); err != nil {
		panic(err)
	}
	defer natsComp.Stop()

	// Create a test event
	evt := pubsub.NewEvent("test-event", nil, nil, map[string]interface{}{
		"message": "Hello NATS!",
		"time":    time.Now().String(),
	})
	evt.SetChannel("test-channel")

	// Subscribe to channel
	ch, close := natsComp.Subscribe(context.Background(), "test-channel", "test-event")
	defer close()

	// Start a goroutine to receive messages
	go func() {
		for event := range ch {
			fmt.Printf("Received event: %s\n", event.String())
			fmt.Printf("Data: %v\n", event.GetData())
		}
	}()

	// Publish event
	fmt.Println("Publishing event...")
	if err := natsComp.Publish(context.Background(), "test-channel", evt); err != nil {
		panic(err)
	}

	// Wait a bit to receive the message
	time.Sleep(time.Second)
}
