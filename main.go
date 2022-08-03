package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MusaSSH/go-mongo-rest-example/handler"
	"github.com/MusaSSH/go-mongo-rest-example/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	client, err := DBInit()
	if err != nil {
		log.Fatal(err)
	}

	uh := handler.NewUserHandler(service.NewUserService(client))

	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
	app.Post("/api/user", uh.PostUser)
	app.Get("/api/user/:objectid", uh.GetUser)
	app.Get("/api/user/", uh.GetUser)
	app.Patch("/api/user/:objectid", uh.UpdateUser)

	app.Listen(":4356")
}

func DBInit() (*mongo.Client, error) {
	uri := os.Getenv("MONGOURI")
	if uri == "" {
		return nil, errors.New("MONGOURI is not defined in environment variables.")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	defer cancel()

	return client, nil
}
