package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := NewRedis(config)
	defer redisClient.RedisClient.Close()

	publisher := NewMessagePublisher(redisClient)

	subscriber := NewMessageConsumer(redisClient)

	go subscriber.ConsumerMessages(ctx, []string{"channel1", "channel2", "channel3", "channel4", "channel5"})

	// Publish messages to multiple channels
	publisher.PublishMessages(ctx, Message{Channel: "channel1", Data: "Hello, Redis!"})
	
	// Publish multiple messages to channel1
	publisher.PublishMessages(ctx, Message{Channel: "channel1", Data: map[string]string{"message": "Hello from channel1 - Message 1"}})
	publisher.PublishMessages(ctx, Message{Channel: "channel1", Data: map[string]string{"message": "Hello from channel1 - Message 2"}})

	publisher.PublishMessages(ctx, Message{Channel: "channel2", Data: map[string]string{"message": "Hello from channel2"}})
	publisher.PublishMessages(ctx, Message{Channel: "channel4", Data: map[string]string{"message": "Hello from channel3"}})
	publisher.PublishMessages(ctx, Message{Channel: "channel4", Data: map[string]string{"message": "Hello from channel4"}})
	publisher.PublishMessages(ctx, Message{Channel: "channel5", Data: map[string]string{"message": "Hello from channel5"}})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}
