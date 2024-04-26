package main

import (
	"context"
	"encoding/json"
	"log"
)

type Message struct {
	Channel string
	Data    interface{}
}

type MessagePublisher struct {
	redisClient Redis
}

func NewMessagePublisher(redisClient Redis) *MessagePublisher {
	return &MessagePublisher{redisClient}
}

func (p *MessagePublisher) PublishMessages(ctx context.Context, message Message) {
	serializedMessage, err := json.Marshal(message.Data)
	if err != nil {
		log.Printf("[%s] Failed to serialize message: %v", message.Channel, err)
		return
	}

	err = p.redisClient.RedisClient.Publish(message.Channel, serializedMessage).Err()
	if err != nil {
		log.Printf("[%s] Failed to publish message: %v", message.Channel, err)
	}
}
