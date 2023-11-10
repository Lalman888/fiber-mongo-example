package models

type Product struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Image       string `json:"image" bson:"image"`
	Description string `json:"description" bson:"description"`
	Price       string `json:"price" bson:"price"`
	MinQuantity int    `json:"minQuantity" bson:"minQuantity"`
	SellerId    string `json:"sellerId" bson:"sellerId"`
}

type CreatePDB struct {
	Name        string `json:"name" bson:"name"`
	Image       string `json:"image" bson:"image"`
	Description string `json:"description" bson:"description"`
	Price       string `json:"price" bson:"price"`
	MinQuantity int    `json:"minQuantity" bson:"minQuantity"`
	SellerId    string `json:"sellerId" bson:"sellerId"`
}

type UpdatePTO struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Image       string `json:"image,omitempty" bson:"image,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Price       string `json:"price,omitempty" bson:"price,omitempty"`
	MinQuantity int    `json:"minQuantity,omitempty" bson:"minQuantity,omitempty"`
	SellerId    string `json:"sellerId,omitempty" bson:"sellerId,omitempty"`
}
