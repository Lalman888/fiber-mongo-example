package router

import (
	"github.com/bmdavis419/fiber-mongo-example/common"
	"github.com/bmdavis419/fiber-mongo-example/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddQueryGroup(app *fiber.App) {
	queryGroup := app.Group("/query")

	queryGroup.Get("/", getQueries)
	queryGroup.Get("/:id", getQuery)
	queryGroup.Post("/", createQuery)
	queryGroup.Delete("/:id", deleteQuery)
}

func getQueries(c *fiber.Ctx) error {
	coll := common.GetDBCollection("query")

	// Find all queries
	queries := make([]models.Query, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Iterate over the cursor
	for cursor.Next(c.Context()) {
		query := models.Query{}
		err := cursor.Decode(&query)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		queries = append(queries, query)
	}

	return c.Status(200).JSON(fiber.Map{"data": queries})
}

func getQuery(c *fiber.Ctx) error {
	coll := common.GetDBCollection("query")

	// Find the query
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	query := &models.Query{}
	filter := bson.M{"_id": objectID}
	err = coll.FindOne(c.Context(), filter).Decode(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": query})
}

type QueryBody struct {
	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Message string `json:"message" bson:"message"`
}

func createQuery(c *fiber.Ctx) error {
	coll := common.GetDBCollection("query")

	// New Product struct
	query := new(QueryBody)

	// Parse body into struct
	err := c.BodyParser(query)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Insert new product
	result, err := coll.InsertOne(c.Context(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return product
	return c.Status(201).JSON(fiber.Map{"data": result})
}

func deleteQuery(c *fiber.Ctx) error {
	coll := common.GetDBCollection("query")

	// Find the query
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	filter := bson.M{"_id": objectID}
	_, err = coll.DeleteOne(c.Context(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"msg": "Query deleted successfully",
	})
}
