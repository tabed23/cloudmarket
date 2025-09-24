package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateProductCollection() error {
	productsCollection := GetCollection("products")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := productsCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Println("Product collection and index created successfully!")
	return nil
}

func CreateCategoryCollection() error {
	categoriesCollection := GetCollection("categories")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := categoriesCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	fmt.Println("Category collection and index created successfully!")

	return nil
}
