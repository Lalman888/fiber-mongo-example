package router

import (
	"github.com/bmdavis419/fiber-mongo-example/common"
	"github.com/bmdavis419/fiber-mongo-example/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProductGroup(app *fiber.App) {
	productGroup := app.Group("/products")

	productGroup.Get("/", getProducts)
	productGroup.Get("/:id", getProduct)
	productGroup.Post("/", createProduct)
	productGroup.Put("/:id", updateProduct)
	productGroup.Delete("/:id", deleteProduct)
}

func getProducts(c *fiber.Ctx) error {
	coll := common.GetDBCollection("products")

	// Find all products
	products := make([]models.Product, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Iterate over the cursor
	for cursor.Next(c.Context()) {
		product := models.Product{}
		err := cursor.Decode(&product)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		products = append(products, product)
	}

	return c.Status(200).JSON(fiber.Map{"data": products})
}

func getProduct(c *fiber.Ctx) error {
	coll := common.GetDBCollection("products")

	// Find the product
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	product := models.Product{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": product})
}

type createPTO struct {
	Name        string  `json:"name" bson:"name"`
	Image       string  `json:"image" bson:"image"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
	MinQuantity int     `json:"minQuantity" bson:"minQuantity"`
	SellerId    string  `json:"sellerId" bson:"sellerId"`
}

func createProduct(c *fiber.Ctx) error {
	// Validate the body
	// p := new(models.Product)
	p := new(createPTO)
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	// Create the product
	coll := common.GetDBCollection("products")
	result, err := coll.InsertOne(c.Context(), p)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create product",
			"message": err.Error(),
		})
	}

	// Return the product
	return c.Status(201).JSON(fiber.Map{
		"result": result,
		"msg":    "Product created successfully",
	})
}

type updatePTO struct {
	Name        string  `json:"name,omitempty" bson:"name,omitempty"`
	Image       string  `json:"image,omitempty" bson:"image,omitempty"`
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Price       float64 `json:"price,omitempty" bson:"price,omitempty"`
	MinQuantity int     `json:"minQuantity,omitempty" bson:"minQuantity,omitempty"`
	SellerId    string  `json:"sellerId,omitempty" bson:"sellerId,omitempty"`
}

func updateProduct(c *fiber.Ctx) error {
	// Validate the body
	// p := new(models.Product)
	p := new(updatePTO)
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	// Get the id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	// Update the product
	coll := common.GetDBCollection("products")
	result, err := coll.UpdateOne(c.Context(), bson.M{"_id": objectID}, bson.M{"$set": p})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update product",
			"message": err.Error(),
		})
	}

	// Return the product
	return c.Status(200).JSON(fiber.Map{
		"result": result,
		"msg":    "Product updated successfully",
	})
}

func deleteProduct(c *fiber.Ctx) error {
	// Get the id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	// Delete the product
	coll := common.GetDBCollection("products")
	result, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete product",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result": result,
		"msg":    "Product deleted successfully",
	})
}
