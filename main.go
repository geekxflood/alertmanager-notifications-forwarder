package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	log.Println("AlertManager Notifications Forwarder is starting...")

	app.Get("*", func(c *fiber.Ctx) error {
		// return the forbidden status code
		return c.Status(403).SendString("Forbidden")
	})

	app.Get("/forward", func(c *fiber.Ctx) error {
		// forward the alert to the intended destination
		// it is also mindfull of the protocol to use (HTTPS or SMTP)
		return c.Status(300).SendString("Forwarding the alert")
	})

	log.Fatal(app.Listen(":3000"))

}
