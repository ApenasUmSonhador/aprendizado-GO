package main

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var contador int32 = 0
var mu sync.Mutex

func incrementa() {
	mu.Lock()
	contador++
	mu.Unlock()
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()
		response := fmt.Sprintf("O valor atual do contador é: %d", contador)
		return c.SendString(response)
	})

	app.Get("/incrementa", func(c *fiber.Ctx) error {
		incrementa()
		return c.SendString("Contador incrementado")
	})

	app.Listen(":3000")
}
