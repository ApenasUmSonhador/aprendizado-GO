package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "./") // Servir todos os arquivos estáticos no diretório atual

	app.Listen(":3000")
}

