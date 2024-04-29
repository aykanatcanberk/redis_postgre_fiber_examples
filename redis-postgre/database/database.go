package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
)

var PgConn *pgx.Conn
var RedisClient *redis.Client

func InitPostgres() {
	var err error
	PgConn, err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/taskdb")
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}
	fmt.Println("Connected to PostgreSQL")

	// Tasks tablosunu oluştur
	_, err = PgConn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			header VARCHAR(255) NOT NULL,
			description TEXT,
			creation_time TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Unable to create tasks table: %v", err)
	}
	fmt.Println("Tasks table created")
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
}

// Diğer veritabanı işlemleri fonksiyonları buraya gelecek
