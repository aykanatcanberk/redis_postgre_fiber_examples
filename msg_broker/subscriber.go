package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

type MessageConsumer struct {
	redisClient  Redis
	subscription *redis.PubSub
}

func NewMessageConsumer(redis Redis) *MessageConsumer {
	return &MessageConsumer{
		redisClient: redis,
	}
}

func (c *MessageConsumer) ConsumerMessages(ctx context.Context, channels []string) {
	for _, channel := range channels {
		go c.handleCustomType1Logic(ctx, channel)
	}
}

func (c *MessageConsumer) handleCustomType1Logic(ctx context.Context, channel string) {
	consumerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("[%s] Consumer started listening...\n", channel)

	c.subscription = c.redisClient.RedisClient.Subscribe(channel)
	defer c.subscription.Close()

	messageChannel := c.subscription.Channel()

	for {
		select {
		case <-consumerCtx.Done():
			log.Printf("[%s] Consumer stopped listening...\n", channel)
			return
		case msg := <-messageChannel:
			var messageData interface{}
			err := json.Unmarshal([]byte(msg.Payload), &messageData)
			if err != nil {
				log.Printf("[%s] Failed to deserialize message: %v", channel, err)
				continue
			}

			fmt.Printf("[%s] Received message: %+v\n", channel, messageData)
		}
	}
}
