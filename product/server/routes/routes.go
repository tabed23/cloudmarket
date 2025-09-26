// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tabed23/cloudmarket-product/server/api"
)

func SetupRoutes(router *gin.Engine, categoryAPI *api.CategoryAPI, productAPI *api.ProductAPI) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "product-service"})
	})
	 
	// API version group
	v1 := router.Group("/api/v1")
	{
		// Category routes
		categories := v1.Group("/categories")
		{
			// Basic CRUD operations
			categories.POST("", categoryAPI.CreateCategoryAPI)
			categories.GET("/:id", categoryAPI.GetCategoryAPI)
			categories.PUT("/:id", categoryAPI.UpdateCategoryAPI)
			categories.DELETE("/:id", categoryAPI.DeleteCategoryAPI)
			
			// List and search operations
			categories.GET("", categoryAPI.ListCategoriesAPI)
			categories.GET("/search", categoryAPI.SearchCategoriesAPI)
			categories.GET("/filter", categoryAPI.FilterCategoriesAPI)
			
			// Category-Product relationship operations
			categories.POST("/assign-product", categoryAPI.AssignProductToCategoryAPI)
			categories.POST("/remove-product", categoryAPI.RemoveProductFromCategoryAPI)
			
			// Products in category
			categories.GET("/:id/products", categoryAPI.ListProductsInCategoryAPI)
			
			// Category hierarchy
			categories.GET("/:id/hierarchy", categoryAPI.GetCategoryHierarchyAPI)
		}
		
		// Product routes
		products := v1.Group("/products")
		{
			// Basic CRUD operations
			products.POST("", productAPI.CreateProductAPI)
			products.GET("/:id", productAPI.GetProductAPI)
			products.PUT("/:id", productAPI.UpdateProductAPI)
			products.DELETE("/:id", productAPI.DeleteProductAPI)
			
			// List, search, and filter operations
			products.GET("", productAPI.ListProductsAPI)
			products.GET("/search", productAPI.SearchProductsAPI)
			products.GET("/filter", productAPI.FilterProductsAPI)
			
			// Categories for product
			products.GET("/:id/categories", categoryAPI.ListCategoriesForProductAPI)
		}
	}
}