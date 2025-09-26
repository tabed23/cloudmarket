package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type MongoConfig struct {
	db *mongo.Client
}

func ConnectMongo(uri string, dbName string)(*MongoConfig, error) {
	client, err :=  mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return  nil, err
	}
	return &MongoConfig{db: client}, nil
}

func (m *MongoConfig)GetCollection(dbname, collectionName string)*mongo.Collection{
	return  m.db.Database(dbname).Collection(collectionName)
}
func (m *MongoConfig)DisconnectMongo()error{
	return m.db.Disconnect(context.TODO())
}


func (m *MongoConfig)CreateProductCollection(dbname string) error {
	productsCollection := m.GetCollection(dbname, "products")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := productsCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoConfig)CreateCategoryCollection(dbName string) error {
	categoriesCollection := m.GetCollection(dbName, "categories")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := categoriesCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}


	return nil
}
