package db

import "log"

func Initialize(uri, dbName string) {
	err := ConnectMongo(uri, dbName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func (){
		err := DisconnectMongo()
		if err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()
	err = CreateProductCollection()
	if err != nil {
		log.Fatal("Failed to create product collection:", err)
	}
	err =  CreateCategoryCollection()
	if err != nil {
		log.Fatal("Failed to create category collection:", err)
	}
	log.Println("MongoDB connected and collections initialized")
}