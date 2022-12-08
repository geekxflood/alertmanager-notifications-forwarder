package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// TODO: do logic is alerts is new or not and handle resolve state
// TODO: add function to send email

func forwardAlert(c *fiber.Ctx, redisClient *redis.Client) {
	var alert AlertManagerPayloadObject
	// unmarshal the alert
	err := json.Unmarshal([]byte(c.Body()), &alert)
	if err != nil {
		log.Println("Error unmarshalling the alert")
		log.Println(err)
		return
	}

	// check if the alert is a recovery
	if alert.Status == "resolved" {
		log.Println("Recovery alert")
		return
	}

	// check if the alert is a new alert
	if alert.Status == "firing" {
		log.Println("New alert")
	}

	// forward the alert to the intended destination
	// it is also mindfull of the protocol to use (HTTPS or SMTP)

}

func main() {

	var portFiber string
	var reedPwd string

	// check is a .env file is present
	// if yes, load it
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")

		portFiber = "3000"
		reedPwd = ""

	} else {
		log.Println(".env file found")
		portFiber = os.Getenv("PORT")
		reedPwd = os.Getenv("REDIS_PASSWORD")
	}

	appFiber := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	log.Println("AlertManager Notifications Forwarder is starting...")
	log.Println("Application listening on port", portFiber)
	log.Println("Connecting to Redis")
	// Instantiate a new Redis client
	reedClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: reedPwd,
		DB:       0,
	})

	// Ping the Redis server to check if the connection is working
	_, errRedis := reedClient.Ping().Result()

	if errRedis != nil {
		log.Println("Error connecting to Redis")
		log.Fatalln(errRedis)
	}

	log.Println("Connected to Redis")

	appFiber.Get("*", func(c *fiber.Ctx) error {
		// return the forbidden status code
		return c.Status(403).SendString("Forbidden")
	})

	appFiber.Get("/forward", func(c *fiber.Ctx) error {

		// forward the alert to the intended destination
		// it is also mindfull of the protocol to use (HTTPS or SMTP)
		go forwardAlert(c, reedClient)
		return c.Status(300).SendString("Forwarding the alert")
	})

	log.Fatalln(appFiber.Listen(":" + portFiber))

}
