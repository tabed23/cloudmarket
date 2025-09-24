package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

func ConnectMongo(uri string, dbName string)error {
	var err error
	client, err =  mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return  err
	}
	err =  client.Connect(context.TODO())
	if err != nil {
		return  err
	}
	database = client.Database(dbName)
	return  nil
}

func GetCollection(collectionName string)*mongo.Collection{
	return database.Collection(collectionName)
}
func DisconnectMongo()error{
	return client.Disconnect(context.TODO())
}