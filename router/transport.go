package router

import (
	"github.com/bmdavis419/fiber-mongo-example/common"
	"github.com/bmdavis419/fiber-mongo-example/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTransportGroup(app *fiber.App) {
	transportGroup := app.Group("/transports")

	transportGroup.Get("/", getTransports)
	transportGroup.Get("/:id", getTransport)
	transportGroup.Post("/", createTransport)
	transportGroup.Put("/:id", updateTransport)
	transportGroup.Delete("/:id", deleteTransport)
}

func getTransports(c *fiber.Ctx) error {
	coll := common.GetDBCollection("transports")

	// Find all transports
	transports := make([]models.Transport, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Iterate over the cursor
	for cursor.Next(c.Context()) {
		transport := models.Transport{}
		err := cursor.Decode(&transport)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		transports = append(transports, transport)
	}

	return c.Status(200).JSON(fiber.Map{"data": transports})
}

func getTransport(c *fiber.Ctx) error {
	coll := common.GetDBCollection("transports")

	// Find the transport
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

	transport := models.Transport{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&transport)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": transport})
}

type TransportQuery struct {
	Name        string   `json:"name" bson:"name"`
	Logo        string   `json:"logo" bson:"logo"`
	Phone       string   `json:"phone" bson:"phone"`
	Sevices     []string `json:"services" bson:"services"`
	Price       float64  `json:"price" bson:"price"`
	MinQuantity int      `json:"minQuantity" bson:"minQuantity"`
	Address     string   `json:"address" bson:"address"`
	Available   bool     `json:"available" bson:"available"`
	Rating      float64  `json:"rating" bson:"rating"`
}

func createTransport(c *fiber.Ctx) error {
	// Validate the body
	// t := new(models.Transport)
	t := new(TransportQuery)
	if err := c.BodyParser(t); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	// Create the transport
	coll := common.GetDBCollection("transports")
	result, err := coll.InsertOne(c.Context(), t)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create transport",
			"message": err.Error(),
		})
	}

	// Return the transport
	return c.Status(201).JSON(fiber.Map{
		"result": result,
	})
}

func updateTransport(c *fiber.Ctx) error {
	// Validate the body
	t := new(models.Transport)
	if err := c.BodyParser(t); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	// Get the ID
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

	// Update the transport
	coll := common.GetDBCollection("transports")
	result, err := coll.UpdateOne(c.Context(), bson.M{"_id": objectID}, bson.M{"$set": t})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update transport",
			"message": err.Error(),
		})
	}

	// Return the transport
	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}

func deleteTransport(c *fiber.Ctx) error {
	// Get the ID
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

	// Delete the transport
	coll := common.GetDBCollection("transports")
	result, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete transport",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}

func AddEnquiryGroup(app *fiber.App) {
	enquiryGroup := app.Group("/enquiries")

	enquiryGroup.Get("/", getEnquiries)
	enquiryGroup.Get("/:id", getEnquiry)
	enquiryGroup.Post("/", createEnquiry)
	enquiryGroup.Put("/:id", updateEnquiry)
	enquiryGroup.Delete("/:id", deleteEnquiry)
}

func getEnquiries(c *fiber.Ctx) error {
	coll := common.GetDBCollection("enquiries")

	// Find all enquiries
	enquiries := make([]models.GenerateEnquiry, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Iterate over the cursor
	for cursor.Next(c.Context()) {
		enquiry := models.GenerateEnquiry{}
		err := cursor.Decode(&enquiry)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		enquiries = append(enquiries, enquiry)
	}

	return c.Status(200).JSON(fiber.Map{"data": enquiries})
}

func getEnquiry(c *fiber.Ctx) error {
	coll := common.GetDBCollection("enquiries")

	// Find the enquiry
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

	enquiry := models.GenerateEnquiry{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&enquiry)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": enquiry})
}

type EnquiryQuery struct {
	TransportId     string `json:"transportId" bson:"transportId"`
	ProductId       string `json:"productId" bson:"productId"`
	Quantity        int    `json:"quantity" bson:"quantity"`
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress"`
	DateOfDelivery  string `json:"dateOfDelivery" bson:"dateOfDelivery"`
	Status          string `json:"status" bson:"status"`
}

func createEnquiry(c *fiber.Ctx) error {
	// Validate the body
	// e := new(models.GenerateEnquiry)
	e := new(EnquiryQuery)
	if err := c.BodyParser(e); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	// Create the enquiry
	coll := common.GetDBCollection("enquiries")
	result, err := coll.InsertOne(c.Context(), e)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create enquiry",
			"message": err.Error(),
		})
	}

	// Return the enquiry
	return c.Status(201).JSON(fiber.Map{
		"result": result,
	})
}

type EnquiryQueryUpdate struct {
	TransportId     string `json:"transportId,omitempty" bson:"transportId,omitempty"`
	ProductId       string `json:"productId,omitempty" bson:"productId,omitempty"`
	Quantity        int    `json:"quantity,omitempty" bson:"quantity,omitempty"`
	DeliveryAddress string `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	DateOfDelivery  string `json:"dateOfDelivery,omitempty" bson:"dateOfDelivery,omitempty"`
	Status          string `json:"status,omitempty" bson:"status,omitempty"`
}

func updateEnquiry(c *fiber.Ctx) error {
	// Validate the body
	e := new(EnquiryQueryUpdate)
	if err := c.BodyParser(e); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	// Get the ID
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

	// Update the enquiry
	coll := common.GetDBCollection("enquiries")
	result, err := coll.UpdateOne(c.Context(), bson.M{"_id": objectID}, bson.M{"$set": e})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update enquiry",
			"message": err.Error(),
		})
	}

	// Return the enquiry
	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}

func deleteEnquiry(c *fiber.Ctx) error {
	// Get the ID
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

	// Delete the enquiry
	coll := common.GetDBCollection("enquiries")
	result, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete enquiry",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}
