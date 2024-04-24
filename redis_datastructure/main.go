package main

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	// Redis'e bağlan
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Redis şifresi varsa buraya yazın, yoksa boş bırakın
		DB:       0,  // Varsayılan Redis veritabanı
	})

	app := fiber.New()

	// Strings (Dizgiler)
	app.Post("/set", setString)
	app.Get("/get/:key", getString)


	// Lists (Listeler)
	app.Post("/lpush", lpush)
	app.Get("/lrange/:key", lrange)

	// Sets (Kümeler)
	app.Post("/sadd", sadd)
	app.Get("/smembers/:key", smembers)

	// Sorted Sets (Sıralı Kümeler)
	app.Post("/zadd", zadd)
	app.Get("/zrange/:key", zrange)

	// Hashes (Hash Tabloları)
	app.Post("/hset", hset)
	app.Get("/hgetall/:key", hgetall)

	// Bitmaps (Bit Eşlemeleri)
	app.Post("/setbit", setbit)
	app.Get("/getbit/:key", getbit)

	// HyperLogLogs
	app.Post("/pfadd", pfadd)
	app.Get("/pfcount/:key", pfcount)

	// Geospatial Indexes (Coğrafi Dizinler)
	app.Post("/geoadd", geoadd)
	app.Get("/geopos/:key", geopos)

	// Fiber uygulamasını 8080 portunda başlat
	log.Fatal(app.Listen(":8080"))
}

// Strings (Dizgiler)
func setString(c *fiber.Ctx) error {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.Set(ctx, req.Key, req.Value, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

// Strings (Dizgiler)
func getString(c *fiber.Ctx) error {
	key := c.Params("key")
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "value": val})
}

// Lists (Listeler)
func lpush(c *fiber.Ctx) error {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.LPush(ctx, req.Key, req.Value).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

func lrange(c *fiber.Ctx) error {
	key := c.Params("key")
	vals, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "values": vals})
}

// Sets (Kümeler)
func sadd(c *fiber.Ctx) error {
	var req struct {
		Key    string `json:"key"`
		Member string `json:"member"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.SAdd(ctx, req.Key, req.Member).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

func smembers(c *fiber.Ctx) error {
	key := c.Params("key")
	members, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "members": members})
}

// Sorted Sets (Sıralı Kümeler)
func zadd(c *fiber.Ctx) error {
	var req struct {
		Key    string  `json:"key"`
		Score  float64 `json:"score"`
		Member string  `json:"member"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.ZAdd(ctx, req.Key, &redis.Z{Score: req.Score, Member: req.Member}).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

func zrange(c *fiber.Ctx) error {
	key := c.Params("key")
	vals, err := rdb.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "values": vals})
}

// Hashes (Hash Tabloları)
func hset(c *fiber.Ctx) error {
	var req struct {
		Key   string `json:"key"`
		Field string `json:"field"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.HSet(ctx, req.Key, req.Field, req.Value).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

func hgetall(c *fiber.Ctx) error {
	key := c.Params("key")
	vals, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "values": vals})
}

// Bitmaps (Bit Eşlemeleri)
func setbit(c *fiber.Ctx) error {
	var req struct {
		Key    string `json:"key"`
		Offset int64  `json:"offset"`
		Value  int    `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.SetBit(ctx, req.Key, req.Offset, req.Value).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

func getbit(c *fiber.Ctx) error {
	key := c.Params("key")
	offset, err := strconv.ParseInt(c.Query("offset", "0"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid offset"})
	}
	val, err := rdb.GetBit(ctx, key, offset).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "offset": offset, "value": val})
}

// HyperLogLogs
func pfadd(c *fiber.Ctx) error {
	var req struct {
		Key      string   `json:"key"`
		Elements []string `json:"elements"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.PFAdd(ctx, req.Key, req.Elements).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

func pfcount(c *fiber.Ctx) error {
	key := c.Params("key")
	count, err := rdb.PFCount(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "count": count})
}

// Geospatial Indexes (Coğrafi Dizinler)
func geoadd(c *fiber.Ctx) error {
	var req struct {
		Key       string  `json:"key"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Member    string  `json:"member"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := rdb.GeoAdd(ctx, req.Key, &redis.GeoLocation{
		Name:      req.Member,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	}).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "OK"})
}

func geopos(c *fiber.Ctx) error {
	key := c.Params("key")
	members := c.Query("members")
	pos, err := rdb.GeoPos(ctx, key, members).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"key": key, "positions": pos})
}
