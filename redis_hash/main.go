package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
)

type Task struct {
	ID           string    `json:"id"`
	Header       string    `json:"header"`
	Description  string    `json:"description"`
	CreationTime time.Time `json:"creation_time"`
}

var redisPool *redis.Pool

func main() {
	// Redis bağlantısı oluştur
	initRedis()

	// Fiber uygulaması oluştur
	app := fiber.New()

	// HTTP endpoint'leri tanımla
	setupRoutes(app)

	// Uygulamayı dinle
	log.Fatal(app.Listen(":3000"))
}

func initRedis() {
	// Redis bağlantı havuzunu oluştur
	redisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func setupRoutes(app *fiber.App) {
	app.Post("/tasks", createTask)
	app.Get("/tasks", listTasks)
	app.Get("/tasks/:id", getTask)
	app.Put("/tasks/:id", updateTask)
	app.Delete("/tasks/:id", deleteTask)
}

func createTask(c *fiber.Ctx) error {
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return err
	}

	// Görevi Redis'e ekle
	conn := redisPool.Get()
	defer conn.Close()

	// Görev ID oluştur
	task.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	// Görevi hash olarak kaydet
	_, err := conn.Do("HMSET", redis.Args{}.Add("task:"+task.ID).
		Add("ID", task.ID).
		Add("Header", task.Header).
		Add("Description", task.Description).
		Add("CreationDate", time.Now().Format("2006-01-02"))...) // Yalnızca gün, ay ve yıl
	if err != nil {
		return err
	}

	return c.JSON(task)
}

func listTasks(c *fiber.Ctx) error {
	conn := redisPool.Get()
	defer conn.Close()

	// Tüm görevleri getir
	keys, err := redis.Strings(conn.Do("KEYS", "task:*"))
	if err != nil {
		return err
	}

	tasks := make([]Task, 0)
	for _, key := range keys {
		taskMap, err := redis.StringMap(conn.Do("HGETALL", key))
		if err != nil {
			return err
		}

		// Görev oluştur ve değerleri atama
		var t Task
		t.ID = taskMap["ID"]
		t.Header = taskMap["Header"]
		t.Description = taskMap["Description"]
		t.CreationTime, _ = time.Parse("2006-01-02", taskMap["CreationDate"]) // Yalnızca gün, ay ve yıl

		tasks = append(tasks, t)
	}

	return c.JSON(tasks)
}

func getTask(c *fiber.Ctx) error {
	id := c.Params("id")

	conn := redisPool.Get()
	defer conn.Close()

	taskMap, err := redis.StringMap(conn.Do("HGETALL", "task:"+id))
	if err != nil {
		return err
	}

	if len(taskMap) == 0 {
		return fiber.ErrNotFound
	}

	// Görev oluştur ve değerleri atama
	var t Task
	t.ID = taskMap["ID"]
	t.Header = taskMap["Header"]
	t.Description = taskMap["Description"]
	t.CreationTime, _ = time.Parse("2006-01-02", taskMap["CreationDate"]) // Yalnızca gün, ay ve yıl

	return c.JSON(t)
}

func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")

	// Yeni görev nesnesi oluştur
	newTask := new(Task)
	if err := c.BodyParser(newTask); err != nil {
		return err
	}

	// Görevi Redis'te güncelle
	conn := redisPool.Get()
	defer conn.Close()

	// Görevin mevcut bilgilerini al
	existingTask, err := getTaskByID(id, conn)
	if err != nil {
		return err
	}

	// Mevcut görevin ID'sini kullanarak güncelleme yap
	_, err = conn.Do("HMSET", redis.Args{}.Add("task:"+id).
		Add("ID", existingTask.ID).
		Add("Header", newTask.Header).
		Add("Description", newTask.Description).
		Add("CreationDate", existingTask.CreationTime.Format("2006-01-02"))...) // Yalnızca gün, ay ve yıl
	if err != nil {
		return err
	}

	// Güncellenmiş görevi döndür
	return c.JSON(existingTask)
}

func getTaskByID(id string, conn redis.Conn) (*Task, error) {
	taskMap, err := redis.StringMap(conn.Do("HGETALL", "task:"+id))
	if err != nil {
		return nil, err
	}

	if len(taskMap) == 0 {
		return nil, fiber.ErrNotFound
	}

	// Görev oluştur ve değerleri atama
	var t Task
	t.ID = taskMap["ID"]
	t.Header = taskMap["Header"]
	t.Description = taskMap["Description"]
	t.CreationTime, _ = time.Parse("2006-01-02", taskMap["CreationDate"]) // Yalnızca gün, ay ve yıl

	return &t, nil
}

func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", "task:"+id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
