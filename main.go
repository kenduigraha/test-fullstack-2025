package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/net/context"
)

var ctx = context.Background()

// Struct user sesuai format Redis
type User struct {
	RealName string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"` // hashed sha1
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Fungsi hash password dengan SHA1
func sha1Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	// Inisialisasi Redis v9
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	app := fiber.New()

	// Define a route for the GET method on the root path '/'
  app.Get("/", func(c fiber.Ctx) error {
      // Send a string response to the client
      return c.SendString("Hello, World 👋!")
  })
	
	app.Post("/login", func(c fiber.Ctx) error {
		var req LoginRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// check key in Redis
		key := fmt.Sprintf("login_%s", req.Username)
		val, err := rdb.Get(ctx, key).Result()

		if err == redis.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Redis error",
			})
		}

		var user User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user data",
			})
		}

		// check password
		if user.Password != sha1Hash(req.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		return c.JSON(fiber.Map{
			"message":  "Login successful",
			"realname": user.RealName,
			"email":    user.Email,
		})
	})

	// Run server
	log.Fatal(app.Listen(":3000"))
}
