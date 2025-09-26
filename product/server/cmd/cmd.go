package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tabed23/cloudmarket-product/server/api"
	"github.com/tabed23/cloudmarket-product/server/db"
	"github.com/tabed23/cloudmarket-product/server/middleware"
	"github.com/tabed23/cloudmarket-product/server/repository"
	"github.com/tabed23/cloudmarket-product/server/repository/store"
	"github.com/tabed23/cloudmarket-product/server/routes"
	"github.com/tabed23/cloudmarket-product/server/service/category"
	"github.com/tabed23/cloudmarket-product/server/service/product"
)

func RunServer() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Environment variables loaded successfully")
	
	uri := getEnv("URI", "mongodb://localhost:27017")
	dbName := getEnv("DATABASE_NAME", "product_service_db")
	port := getEnv("PORT", "8080")
	
	db, err := db.ConnectMongo(uri, dbName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	db.CreateCategoryCollection(dbName)
	db.CreateProductCollection(dbName)

	if err := middleware.InitLogDirectories(); err != nil {
		log.Printf("Warning: Failed to create log directories: %v", err)
	}
	// Initialize application
	app := initializeApplication(db, dbName)

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorLoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())
	
	// Setup routes
	routes.SetupRoutes(router, app.CategoryAPI, app.ProductAPI)

	// Start server
	fmt.Printf("Starting server on port :%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type Application struct {
	CategoryRepo repository.CategoryRepository
	ProductRepo  repository.ProductRepository
	CategorySvc  repository.CategoryRepository  // Your service returns repository interface
	ProductSvc   repository.ProductRepository   // Your service returns repository interface
	CategoryAPI  *api.CategoryAPI
	ProductAPI   *api.ProductAPI
}

func initializeApplication(db *db.MongoConfig, dbName string) *Application {
	// Initialize collections using db.GetCollection
	categoryCollection := db.GetCollection(dbName, "categories")
	productCollection := db.GetCollection(dbName,"products")

	// Initialize repositories
	categoryRepo := store.NewCategoryStore(categoryCollection, productCollection)
	productRepo := store.NewProductStore(productCollection)

	// Initialize services (they return repository interfaces based on your code)
	categorySvc := category.NewCategoryService(categoryRepo)
	productSvc := product.NewProductService(productRepo)

	// Initialize APIs - Fix the pointer issue
	// CategoryAPI expects category.CategoryService, but your service returns repository.CategoryRepository
	// You need to cast or fix the service layer
	categoryServiceImpl := categorySvc.(*category.CategoryService) // Cast back to concrete type
	productServiceImpl := productSvc.(*product.ProductService)     // Cast back to concrete type

	categoryAPI := api.NewCategoryAPI(*categoryServiceImpl)
	productAPI := api.NewProductAPI(productServiceImpl)

	return &Application{
		CategoryRepo: categoryRepo,
		ProductRepo:  productRepo,
		CategorySvc:  categorySvc,
		ProductSvc:   productSvc,
		CategoryAPI:  categoryAPI,
		ProductAPI:   productAPI,
	}
}