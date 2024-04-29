package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"task/database"
	"task/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Post("/task", createTask)
	app.Get("/task/:id", getTask)
	app.Put("/task/:id", updateTask)
	app.Delete("/task/:id", deleteTask)
	app.Get("/redis/keys", listRedisKeys)
}

func createTask(c *fiber.Ctx) error {
	task := new(models.Task)

	if err := c.BodyParser(task); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	task.CreationTime = time.Now()

	_, err := database.PgConn.Exec(context.Background(), "INSERT INTO tasks (header, description, creation_time) VALUES ($1, $2, $3)", task.Header, task.Description, task.CreationTime)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(task)
}

func getTask(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	// Öncelikle Redis'te arayalım
	val, err := database.RedisClient.Get(ctx, id).Result()
	if err == nil {
		// Cache'te bulundu
		task := new(models.Task)
		if err := json.Unmarshal([]byte(val), task); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(task)
	}

	// Redis'te bulunamazsa, PostgreSQL'den arayalım
	task := new(models.Task)
	err = database.PgConn.QueryRow(ctx, "SELECT id, header, description, creation_time FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Header, &task.Description, &task.CreationTime)
	if err != nil {
		return c.Status(404).SendString("Task not found")
	}

	// Sonucu Redis'e kaydedelim
	jsonData, err := json.Marshal(task)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	database.RedisClient.Set(ctx, fmt.Sprintf("%d", task.ID), string(jsonData), 30*time.Minute) // 30 dakika boyunca cache'le

	return c.JSON(task) // Fonksiyonun sonunda başarıyla Task'ı JSON olarak döndür
}

func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task := new(models.Task)

	if err := c.BodyParser(task); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Veritabanında güncelleme işlemi yapılacak
	ctx := context.Background()
	_, err := database.PgConn.Exec(ctx, "UPDATE tasks SET header = $1, description = $2 WHERE id = $3", task.Header, task.Description, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// Güncelleme başarılı oldu, HTTP 200 OK dön
	return c.SendStatus(fiber.StatusOK)
}

func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	// Veritabanında silme işlemi yapılacak
	ctx := context.Background()
	_, err := database.PgConn.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// Silme başarılı oldu, HTTP 200 OK dön
	return c.SendStatus(fiber.StatusOK)
}

func listRedisKeys(c *fiber.Ctx) error {
	ctx := context.Background()

	// Tüm anahtarları al
	keys, err := database.RedisClient.Keys(ctx, "*").Result()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// Anahtarlarla ilişkili değerleri al
	var results []map[string]interface{}
	for _, key := range keys {
		val, err := database.RedisClient.Get(ctx, key).Result()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		data := make(map[string]interface{})
		data[key] = val
		results = append(results, data)
	}

	return c.JSON(results)
}
