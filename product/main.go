package main

import (
	"fmt"

	"github.com/tabed23/cloudmarket-product/server/db"
)


func main(){
	// MongoDB URI and database name
    uri := "mongodb://admin:admin123@localhost:27017"
    dbName := "product_service_db"

    db.Initialize(uri, dbName)

    fmt.Println("Database initialized successfully!")
}