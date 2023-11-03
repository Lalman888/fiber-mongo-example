package models

type Product struct {
	ID          string  `json:"_id" bson:"_id"`
	Name        string  `json:"name" bson:"name"`
	Image       string  `json:"image" bson:"image"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
	MinQuantity int     `json:"minQuantity" bson:"minQuantity"`
	SellerId    string  `json:"sellerId" bson:"sellerId"`
}
