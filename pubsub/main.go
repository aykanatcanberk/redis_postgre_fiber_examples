package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
)

func main() {
	// Fiber uygulamasını başlat
	app := fiber.New()

	// Redis istemcisini oluştur
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Docker Compose ile ağ içerisindeki Redis'e erişim sağlamak için
	})

	// Mesaj kuyruğu ismi
	queueName := "message_queue"

	// Kanal oluştur
	msgCh := make(chan string)

	// WaitGroup oluştur
	var wg sync.WaitGroup

	// Abone olacağımız kanalı belirt
	wg.Add(1)
	go subscriber(rdb, queueName, msgCh, &wg)

	// POST endpoint'i oluştur
	app.Post("/publish", func(c *fiber.Ctx) error {
		var requestBody struct {
			Message string `json:"message"`
		}

		// JSON body'yi oku
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid JSON format")
		}

		// Mesajı yayınla ve kuyruğa ekle
		err := publishMessage(c, rdb, queueName, requestBody.Message)
		if err != nil {
			log.Println("Error publishing message:", err)
			return c.Status(http.StatusInternalServerError).SendString("Error publishing message")
		}

		return c.SendString("Message queued: " + requestBody.Message)
	})

	// Kanaldan gelen mesajları dinleme
	go func() {
		for {
			select {
			case msg := <-msgCh:
				// Yeni bir mesaj alındığında bir goroutine başlat ve mesajı işle
				go func(msg string) {
					fmt.Println("Received message from queue:", msg)
				}(msg)
			}
		}
	}()

	// Fiber uygulamasını dinlemeye başla
	log.Fatal(app.Listen(":3000"))

	// WaitGroup'i bekleyerek programın sonlanmasını sağla
	wg.Wait()
}

// Mesaj yayınlama ve kuyruğa ekleme fonksiyonu
func publishMessage(c *fiber.Ctx, rdb *redis.Client, queueName, message string) error {
	// Mesajı Redis kanalına yayınla
	err := rdb.Publish(c.Context(), "messages", message).Err()
	if err != nil {
		return err
	}

	// Mesajı Redis kuyruğuna ekle
	_, err = rdb.LPush(c.Context(), queueName, message).Result()
	if err != nil {
		return err
	}

	return nil
}

// Abone olma ve kuyruktan mesajları alma fonksiyonu
func subscriber(rdb *redis.Client, queueName string, msgCh chan<- string, wg *sync.WaitGroup) {
	// WaitGroup'i azaltarak işlem tamamlandığında haber ver
	defer wg.Done()

	// context.Background() ile bir bağlam oluştur
	ctx := context.Background()

	// Sonsuz bir döngü içinde mesajları dinle
	for {
		// Mesaj kuyruğundan mesaj al
		msg, err := rdb.BRPopLPush(ctx, queueName, queueName, 0).Result()
		if err != nil {
			log.Println("Error receiving message from queue:", err)
			continue
		}
		// Kanala mesajı gönder
		msgCh <- msg
	}
}
