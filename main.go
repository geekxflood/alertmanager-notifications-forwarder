package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// emailSender sends an email to the target user
func emailSender() {
	// Create a new TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Connect to the email server
	conn, err := tls.Dial("tcp", "mail.example.com:465", tlsConfig)
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	// Set up authentication information
	auth := smtp.PlainAuth("", "user@example.com", "password", "mail.example.com")

	// Set up the email message
	to := []string{"recipient@example.com"}
	msg := []byte("To: recipient@example.com\r\n" +
		"Subject: Hello from Golang\r\n" +
		"\r\n" +
		"This is a test email from Golang.\r\n")

	// Send the email
	err = smtp.SendMail("mail.example.com:465", auth, "user@example.com", to, msg)
	if err != nil {
		log.Panic(err)
	}

}

// alertResolvedCheckings checks if the alert is already resolved or not
// if the alert is already resolved, it will return false
func alertChecking(a AlertObject, state bool) (bool, error) {
	// Create a new context for the Redis client
	ctx := context.Background()

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})

	if state {
		// There is an alert firing
		// Check if the alert is already in the database
		// If the alert is already in the database, do nothing
		_, err := rdb.Get(ctx, a.Fingerprint).Result()
		if err != nil {
			// The alert is not in the database
			// Add the alert to the database
			_, err := rdb.Set(ctx, a.Fingerprint, true, 0).Result()
			if err != nil {
				return false, err
			}
			return true, nil
		}

	} else {
		// The alert is resolved
		// Check if the alert is already in the database
		// If the alert is already in the database, delete it
		// If the alert is not in the database, do nothing
		_, err := rdb.Del(ctx, a.Fingerprint).Result()
		if err != nil {
			return false, err
		}

	}
	return true, nil
}

func main() {

	// Check if there is a .env file in the current directory
	// If there is a .env file, load the environment variables from the .env file
	if _, err := os.Stat(".env"); err == nil {
		// There is a .env file in the current directory
		// Load the environment variables from the .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		// There is no .env file in the current directory
		// Set default environment variables
		os.Setenv("REDIS_HOST", "localhost")
		os.Setenv("REDIS_PORT", "6379")
		os.Setenv("APP_PORT", "9847")
	}

	// Create a new Fiber app
	// Disable the startup message
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)

	log.Println("Starting the server on port", os.Getenv("APP_PORT"))

	// Create the POST /alert endpoint
	app.Post("/alert", func(c *fiber.Ctx) error {
		// Parse the request body as an Alert object
		// If the request body is not an Alert object, return an error
		var alertBody AlertManagerPayloadObject
		if err := c.BodyParser(&alertBody); err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		// Iterate through the alerts
		for _, alert := range alertBody.Alert {
			// Check if the alert is firing or resolved
			if alert.Status == "firing" {
				// The alert is firing
				// Check if the alert is already in the database
				// If the alert is already in the database, do nothing
				// newAlert will be false if the alert is already in the database
				newAlert, err := alertChecking(alert, true)

				if err != nil {
					log.Fatalln(err)
				}
				if newAlert {
					log.Println("New Alert", alert.Labels.Alertname, "is firing")

				}
				// If the alert is resolved, delete it from the database if it is in the database
			} else if alert.Status == "resolved" {
				log.Println("Alert", alert.Labels.Alertname, "is resolved")
				// The alert is resolved, no need to have the `newAlert` variable
				_, err := alertChecking(alert, false)
				if err != nil {
					log.Fatalln(err)
				}

			}
		}
		return c.SendString("Success")
	})

	// Create the GET * endpoint
	// If the user tries to access any other endpoint, return a 403 Forbidden error
	app.Get("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	})

	// Start the server on port as default 9847
	log.Fatalln(app.Listen(":" + os.Getenv("APP_PORT")))

}
