// main.go

package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	app         *fiber.App
)

func main() {
	// Redis client oluştur
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis sunucu adresi
		Password: "",               // Eğer bir şifre kullanıyorsanız buraya ekleyin
		DB:       0,                // Redis veritabanı seçimi
	})

	// Fiber uygulamasını oluştur
	app = fiber.New()

	// Kanal oluşturma endpoint'i
	app.Post("/channel", createChannel)

	// Kanalları listeleme endpoint'i
	app.Get("/channels", listChannels)

	// Kanal silme endpoint'i
	app.Delete("/channel/:name", deleteChannel)

	// Mesaj yayınlama endpoint'i
	app.Post("/publish/:channel", publishMessage)

	// Abone olma endpoint'i
	app.Post("/subscribe/:channel", subscribeChannel)

	// Aboneleri listeleme endpoint'i
	app.Get("/subscribers/:channel", listSubscribers)

	// Kanaldaki mesajları listeleme endpoint'i
	app.Get("/messages/:channel", listMessages)

	// Fiber sunucuyu başlat
	log.Fatal(app.Listen(":3000"))
}

// Kanal oluşturma handler fonksiyonu
func createChannel(c *fiber.Ctx) error {
	channelName := c.Params("name")
	if channelName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Channel name cannot be empty",
		})
	}

	// Redis SUBSCRIBE komutunu kullanarak kanalı oluştur
	pubsub := redisClient.Subscribe(c.Context(), channelName)
	defer pubsub.Close() // Aboneliği kapat

	// Abonelik işlemi sırasında oluşabilecek hataları kontrol et
	if _, err := pubsub.Receive(c.Context()); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Channel %s created", channelName),
	})
}



// Kanalları listeleme handler fonksiyonu
func listChannels(c *fiber.Ctx) error {
	var channels []string

	// SCAN komutu ile anahtarları tara
	iter := redisClient.Scan(c.Context(), 0, "*", 0).Iterator()
	for iter.Next(c.Context()) {
		channels = append(channels, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"channels": channels,
	})
}

// Kanal silme handler fonksiyonu
func deleteChannel(c *fiber.Ctx) error {
	channelName := c.Params("name")
	err := redisClient.Del(c.Context(), channelName).Err()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Channel %s deleted", channelName),
	})
}

// Mesaj yayınlama handler fonksiyonu
func publishMessage(c *fiber.Ctx) error {
	channelName := c.Params("channel")
	message := c.FormValue("message")
	err := redisClient.Publish(c.Context(), channelName, message).Err()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Message published to channel %s", channelName),
	})
}

// Abone olma handler fonksiyonu
func subscribeChannel(c *fiber.Ctx) error {
	channelName := c.Params("channel")
	subscriberName := c.Params("name") // API isteğinden abone olan kişinin adını al
	if subscriberName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Subscriber name cannot be empty",
		})
	}

	// Redis SADD komutunu kullanarak abone ol
	err := redisClient.SAdd(c.Context(), channelName+":subscribers", subscriberName).Err()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Subscribed %s to channel %s", subscriberName, channelName),
	})
}



// Aboneleri listeleme handler fonksiyonu
func listSubscribers(c *fiber.Ctx) error {
	channelName := c.Params("channel")
	subscribers, err := redisClient.SMembers(c.Context(), channelName+":subscribers").Result()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"subscribers": subscribers,
	})
}

// Kanaldaki mesajları listeleme handler fonksiyonu
func listMessages(c *fiber.Ctx) error {
	channelName := c.Params("channel")
	messages, err := redisClient.LRange(c.Context(), channelName+":messages", 0, -1).Result()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"messages": messages,
	})
}
