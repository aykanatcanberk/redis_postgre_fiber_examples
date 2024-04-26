package main

import "github.com/go-redis/redis"

type Redis struct {
	RedisClient *redis.Client
}

func NewRedis(config Config) Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       0,
	})

	return Redis{
		RedisClient: client,
	}
}
