package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	Name        string                `form:"name" bson:"name"`
	Image       *multipart.FileHeader `form:"image" bson:"image"`
	Description string                `form:"description" bson:"description"`
	Price       string                `form:"price" bson:"price"`
	MinQuantity int                   `form:"minQuantity" bson:"minQuantity"`
	SellerId    string                `form:"sellerId" bson:"sellerId"`
}

func createProduct(c *fiber.Ctx) error {
	// Validate the body
	p := new(createPTO)
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Retrieve AWS region from environment variable
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Println("AWS_REGION environment variable is not set.")
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to load AWS S3 config",
			"message": "AWS_REGION environment variable is not set.",
		})
	}

	// Setup S3 uploader
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
	)
	if err != nil {
		log.Printf("error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to load AWS S3 config",
			"message": err.Error(),
		})
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	// Handle product image upload
	imageURL, err := handleProductUpload(c, uploader)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to upload product image",
			"message": err.Error(),
		})
	}

	// Set the imageURL in the product struct

	newData := &models.CreatePDB{
		Name:        p.Name,
		Image:       imageURL,
		Description: p.Description,
		Price:       p.Price,
		MinQuantity: p.MinQuantity,
		SellerId:    p.SellerId,
	}

	// Create the product
	coll := common.GetDBCollection("products")
	result, err := coll.InsertOne(c.Context(), newData)
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

type updatePTODB struct {
	Name        string                `form:"name,omitempty" bson:"name,omitempty"`
	Image       *multipart.FileHeader `form:"image,omitempty" bson:"image,omitempty"`
	Description string                `form:"description,omitempty" bson:"description,omitempty"`
	Price       string                `form:"price,omitempty" bson:"price,omitempty"`
	MinQuantity int                   `form:"minQuantity,omitempty" bson:"minQuantity,omitempty"`
	SellerId    string                `form:"sellerId,omitempty" bson:"sellerId,omitempty"`
}

func updateProduct(c *fiber.Ctx) error {
	// Validate the body
	p := new(models.UpdatePTO)
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

func handleProductUpload(c *fiber.Ctx, uploader *manager.Uploader) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", err
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}

	// Determine Content-Type based on the file extension (you can enhance this logic)
	contentType := determineContentType(file.Filename)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String("grain"),
		Key:                aws.String("grains/gi" + file.Filename),
		Body:               f,
		ACL:                "public-read",
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("inline"), // Set to "inline" to display in the browser

	}, func(u *manager.Uploader) {
		u.PartSize = 6 * 1024 * 1024 // Override the PartSize to 6 MiB
	})

	if err != nil {
		var mu manager.MultiUploadFailure
		if errors.As(err, &mu) {
			fmt.Println("Error:", mu)
			uploadID := mu.UploadID()
			return uploadID, err
		} else {
			fmt.Println("Error:", err.Error())
			return "", err
		}
	}

	return result.Location, nil
}

func determineContentType(filename string) string {
	// You can implement more sophisticated logic to determine the Content-Type
	// For simplicity, this example uses a basic mapping based on file extension
	switch filepath.Ext(filename) {
	case ".pdf":
		return "application/pdf"
	case ".png":
		return "image/png"
	default:
		// Set a default Content-Type or handle unknown types accordingly
		return "application/octet-stream"
	}
}
