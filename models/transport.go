package models

type Transport struct {
	ID          string   `json:"_id" bson:"_id"`
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

type GenerateEnquiry struct {
	ID              string `json:"id" bson:"_id"`
	TransportId     string `json:"transportId" bson:"transportId"`
	ProductId       string `json:"productId" bson:"productId"`
	Quantity        int    `json:"quantity" bson:"quantity"`
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress"`
	DateOfDelivery  string `json:"dateOfDelivery" bson:"dateOfDelivery"`
	Status          string `json:"status" bson:"status"`
}
