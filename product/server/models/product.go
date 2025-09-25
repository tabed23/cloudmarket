package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Brand        string             `bson:"brand"`
	Description  string             `bson:"description"`
	Price        float64            `bson:"price"`
	Discount     string             `bson:"discount"`
	Rating       float64            `bson:"rating"`
	ReviewsCount int                `bson:"reviews_count"`
	Sizes        []string           `bson:"sizes"`
	Colors       []string           `bson:"colors"`
	Images       []string           `bson:"images"`
	Dimensions   string             `bson:"dimensions"`
	CategoryID   primitive.ObjectID `bson:"category_id"`
	Categories   []primitive.ObjectID `bson:"categories,omitempty"` 
	CreatedAt    primitive.DateTime `bson:"created_at"`
}
