package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var rdb *redis.Client

// Task struct
type Task struct {
	ID          int       `json:"id"`
	Header      string    `json:"header"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	// Redis client configuration
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server adresi
		Password: "",               // Redis şifresi
		DB:       0,                // Redis veritabanı indeksi
	})

	// Fiber web framework'u oluşturma
	app := fiber.New()

	// LPUSH: Listenin sol tarafına bir veya birden fazla öğe eklemek için endpoint
	app.Post("/task/lpush", lpushTask)

	// RPUSH: Listenin sağ tarafına bir veya birden fazla öğe eklemek için endpoint
	app.Post("/task/rpush", rpushTask)

	// LPOP: Listenin sol tarafındaki bir öğeyi kaldırmak için endpoint
	app.Get("/task/lpop", lpopTask)

	// RPOP: Listenin sağ tarafındaki bir öğeyi kaldırmak için endpoint
	app.Get("/task/rpop", rpopTask)

	// LRANGE: Belirtilen aralıktaki öğeleri döndürmek için endpoint
	app.Get("/task/lrange", lrangeTask)

	// LINDEX: Belirtilen dizindeki öğeyi döndürmek için endpoint
	app.Get("/task/lindex/:index", lindexTask)

	// LINSERT: Belirli bir öğeyi belirli bir öğeden önce veya sonra eklemek için endpoint
	app.Post("/task/linsert", linsertTask)

	// LLEN: Listenin uzunluğunu döndürmek için endpoint
	app.Get("/task/llen", llenTask)

	// LREM: Belirtilen değeri belirli sayıda listeden kaldırmak için endpoint
	app.Post("/task/lrem", lremTask)

	// LSET: Belirli bir dizindeki öğeyi değiştirmek için endpoint
	app.Post("/task/lset/:index", lsetTask)

	// LTRIM: Listeyi belirtilen aralıkta kırpma için endpoint
	app.Post("/task/ltrim", ltrimTask)

	// Fiber sunucusunu başlatma
	log.Fatal(app.Listen(":3000"))
}

// lpushTask, POST isteklerini dinler ve görevi listenin soluna ekler
func lpushTask(c *fiber.Ctx) error {
	taskID, err := rdb.LPush(c.Context(), "tasks", "Yeni görev").Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Görev Kimliği: %d", taskID))
}

// rpushTask, POST isteklerini dinler ve görevi listenin sağına ekler
func rpushTask(c *fiber.Ctx) error {
	var requestData struct {
		Item string `json:"item"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	taskID, err := rdb.RPush(c.Context(), "tasks", requestData.Item).Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Görev Kimliği: %d", taskID))
}

// lpopTask, GET isteklerini dinler ve listenin solundan görevi çıkarır
func lpopTask(c *fiber.Ctx) error {
	task, err := rdb.LPop(c.Context(), "tasks").Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Çıkarılan Görev: %s", task))
}

// rpopTask, GET isteklerini dinler ve listenin sağından görevi çıkarır
func rpopTask(c *fiber.Ctx) error {
	task, err := rdb.RPop(c.Context(), "tasks").Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Çıkarılan Görev: %s", task))
}

// lrangeTask, GET isteklerini dinler ve belirtilen aralıktaki görevleri döndürür
func lrangeTask(c *fiber.Ctx) error {
	tasks, err := rdb.LRange(c.Context(), "tasks", 0, -1).Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(tasks)
}

// lindexTask, GET isteklerini dinler ve belirtilen dizindeki görevi döndürür
func lindexTask(c *fiber.Ctx) error {
	index, err := strconv.ParseInt(c.Params("index"), 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	task, err := rdb.LIndex(c.Context(), "tasks", index).Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Index %d'deki Görev: %s", index, task))
}

// linsertTask, POST isteklerini dinler ve yeni bir öğeyi var olan bir öğenin önüne veya arkasına ekler
func linsertTask(c *fiber.Ctx) error {
	var requestData struct {
		ExistingItem string `json:"existing_item"`
		NewItem      string `json:"new_item"`
		Position     string `json:"position"` // ÖNCE veya SONRA
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	taskID, err := rdb.LInsert(c.Context(), "tasks", requestData.Position, requestData.ExistingItem, requestData.NewItem).Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Görev Kimliği: %d", taskID))
}

// llenTask, GET isteklerini dinler ve listenin uzunluğunu döndürür
func llenTask(c *fiber.Ctx) error {
	length, err := rdb.LLen(c.Context(), "tasks").Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Liste Uzunluğu: %d", length))
}

// lremTask, POST isteklerini dinler ve belirtilen değeri belirli sayıda listeden kaldırır
func lremTask(c *fiber.Ctx) error {
	removedCount, err := rdb.LRem(c.Context(), "tasks", 1, "Yeni görev").Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fmt.Sprintf("Kaldırılan Görev Sayısı: %d", removedCount))
}

// lsetTask, POST isteklerini dinler ve belirli bir dizindeki öğenin değerini ayarlar
func lsetTask(c *fiber.Ctx) error {
	index, err := strconv.ParseInt(c.Params("index"), 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = rdb.LSet(c.Context(), "tasks", index, "Güncellenmiş görev").Err()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString("Görev başarıyla güncellendi")
}

// ltrimTask, POST isteklerini dinler ve listeyi belirtilen aralıkta kırp
func ltrimTask(c *fiber.Ctx) error {
	err := rdb.LTrim(c.Context(), "tasks", 0, 2).Err()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString("Liste başarıyla kırpıldı")
}
