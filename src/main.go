package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func sendEmail(alert AlertObject) error {
	// Send the email
	return nil
}

func alertFiringChecking(a AlertObject) (bool, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     string(os.Getenv("REDIS_HOST")),
		Password: "",
		DB:       0,
	})

	val, err := rdb.Get(ctx, string(a.Fingerprint)).Result()
	switch {
	case err == redis.Nil:
		log.Println("key does not exist")
		log.Println("Creating new key", string(a.Fingerprint), "with value 1")
		err := rdb.Set(ctx, string(a.Fingerprint), "1", 0).Err()
		if err != nil {
			log.Fatalln(err)
		}
		return true, nil

	case err != nil:
		log.Fatalln("Get failed", err)
		return true, err
	case val == "":
		log.Fatalln("value is empty")
		return true, err
	}

	if val == "1" {
		log.Println("Alert already present", a.Labels.Alertname, "is firing")
		return false, nil
	}

	return true, nil
}

func main() {
	app := fiber.New()

	// Use the logger and recovery middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Create the POST /alert endpoint
	app.Post("/alert", func(c *fiber.Ctx) error {
		// Parse the request body as an Alert object
		var alertBody AlertManagerPayloadObject
		if err := c.BodyParser(&alertBody); err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		for _, alert := range alertBody.Alert {
			if alert.Status == "firing" {

				newAlert, err := alertFiringChecking(alert)

				if err != nil {
					log.Fatalln(err)
				}
				if newAlert {
					log.Println("New Alert", alert.Labels.Alertname, "is firing")

					// New alert is firing, create the report document
					sendEmail(alert)
				}

			} else if alert.Status == "resolved" {
				log.Println("Alert", alert.Labels.Alertname, "is resolved")
				_, err := alertResolvedCheckings(alert)
				if err != nil {
					log.Fatalln(err)
				}

			}
		}
		return c.SendString("Success")
	})

	app.Get("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	})

	// Start the server on port 9847
	log.Fatalln(app.Listen(":9847"))

}
