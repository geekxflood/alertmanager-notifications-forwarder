package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func templater(a AlertObject) (string, error) {

	// Create a new buffer to store the template result
	buff := new(bytes.Buffer)

	// Parse the template file
	t, err := template.ParseFiles("template.html")
	if err != nil {
		return "", err
	}

	// Generate the value from the template
	err = t.ExecuteTemplate(buff, "template.html", a)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

// emailSender sends an email to the target user
func emailSender(a AlertObject) {

	// Generate the email body message
	messageBody, err := templater(a)
	if err != nil {
		log.Println("Error while generating the email template")
		log.Println(err)
		return
	}

	// timeNow := time.Now()
	portVal, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	// Create a new message
	msg := gomail.NewMessage()
	msg.SetHeader("From", os.Getenv("SMTP_FROM"))
	msg.SetHeader("To", os.Getenv("SMTP_TO"))
	msg.SetHeader("Subject", "Alarme "+a.Labels.Severity+" - "+a.Labels.Alertname)
	// Use the file content as the email body and have it generated html from it's template
	msg.SetBody("text/html", messageBody)

	n := gomail.NewDialer(os.Getenv("SMTP_HOST"), portVal, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	n.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		log.Println("Error while sending the email")
		log.Println(err)
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

	// ALERT FIRING
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

			// Send an email
			emailSender(a)
			return true, nil
		}
		// ALERT RESOLVED
	} else {
		// The alert is resolved
		// Check if the alert is already in the database
		// If the alert is already in the database, delete it
		// If the alert is not in the database, do nothing
		_, err := rdb.Del(ctx, a.Fingerprint).Result()
		if err != nil {
			return false, err
		}

		// Send an email
		emailSender(a)

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
		os.Setenv("SMTP_HOST", "localhost")
		os.Setenv("SMTP_PORT", "587")
		os.Setenv("SMTP_USERNAME", "username")
		os.Setenv("SMTP_PASSWORD", "password")
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
