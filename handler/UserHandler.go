package handler

import (
	"net/http"

	"github.com/MusaSSH/go-mongo-rest-example/models"
	"github.com/MusaSSH/go-mongo-rest-example/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	service service.UserService
}

func (h UserHandler) PostUser(c *fiber.Ctx) error {
	var user models.User
	err := c.BodyParser(&user)
	if err != nil {
		return err
	}

	if len(user.Name) < 4 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "length of the name can't be less than 4",
		})
	}

	if len(user.Password) < 8 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "length of the password can't be less than 8",
		})
	}

	oid, err := h.service.Insert(user)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"ObjectId": oid.Hex(),
	})
}

func (h UserHandler) GetUser(c *fiber.Ctx) error {
	var filter models.User
	objectidhex := c.Params("objectid")
	if objectidhex != "" {
		objectid, err := primitive.ObjectIDFromHex(objectidhex)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "not a real objectid",
			})
		}

		filter.Id = objectid
	} else {
		err := c.QueryParser(&filter)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "query parse error",
			})
		}
	}

	user, err := h.service.Get(filter)
	if err == mongo.ErrNoDocuments {
		return c.SendStatus(http.StatusNotFound)
	} else if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusFound).JSON(user)
}

func (h UserHandler) UpdateUser(c *fiber.Ctx) error {
	objectidhex := c.Params("objectid")
	objectid, err := primitive.ObjectIDFromHex(objectidhex)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "not a real objectid",
		})
	}

	var update models.User
	err = c.BodyParser(&update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "query parse error",
		})
	}

	err = h.service.Update(objectid, update)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func NewUserHandler(service service.UserService) UserHandler {
	return UserHandler{
		service: service,
	}
}
