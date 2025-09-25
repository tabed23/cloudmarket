package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/service/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductAPI struct {
	svc *product.ProductService
}

func NewProductAPI(svc *product.ProductService) *ProductAPI {
	return &ProductAPI{svc: svc}
}

// CreateProductAPI creates a new product
func (papi *ProductAPI) CreateProductAPI(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProduct, err := papi.svc.CreateProduct(c.Request.Context(), &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

// GetProductAPI retrieves a single product by ID
func (papi *ProductAPI) GetProductAPI(c *gin.Context) {
	productIDStr := c.Param("id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id format"})
		return
	}

	product, err := papi.svc.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	if product.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProductAPI updates an existing product
func (papi *ProductAPI) UpdateProductAPI(c *gin.Context) {
	productIDStr := c.Param("id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id format"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent changing the ID
	product.ID = productID

	updatedProduct, err := papi.svc.UpdateProduct(c.Request.Context(), productID, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProductAPI deletes a product by ID
func (papi *ProductAPI) DeleteProductAPI(c *gin.Context) {
	productIDStr := c.Param("id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id format"})
		return
	}

	err = papi.svc.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// ListProductsAPI retrieves all products with pagination and optional filtering
func (papi *ProductAPI) ListProductsAPI(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Parse optional filters
	categoryID := c.Query("category_id")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	status := c.Query("status")
	brand := c.Query("brand")

	var products []models.Product
	var err error

	// Build filters if any are provided
	if categoryID != "" || minPrice != "" || maxPrice != "" || status != "" || brand != "" {
		filters := bson.M{}

		if categoryID != "" {
			catID, err := primitive.ObjectIDFromHex(categoryID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id format"})
				return
			}
			filters["$or"] = []bson.M{
				{"category_id": catID},
				{"categories": catID},
			}
		}

		if minPrice != "" || maxPrice != "" {
			priceFilter := bson.M{}
			if minPrice != "" {
				if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
					priceFilter["$gte"] = min
				}
			}
			if maxPrice != "" {
				if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
					priceFilter["$lte"] = max
				}
			}
			if len(priceFilter) > 0 {
				filters["price"] = priceFilter
			}
		}

		if status != "" {
			filters["status"] = status
		}

		if brand != "" {
			filters["brand"] = primitive.Regex{Pattern: brand, Options: "i"}
		}

		products, err = papi.svc.FilterProducts(c.Request.Context(), filters)
	} else {
		products, err = papi.svc.GetAllProducts(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply pagination
	total := len(products)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedProducts := products[start:end]

	c.JSON(http.StatusOK, gin.H{
		"products": paginatedProducts,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// SearchProductsAPI searches products by query string
func (papi *ProductAPI) SearchProductsAPI(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := papi.svc.SearchProducts(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply pagination
	total := len(products)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedProducts := products[start:end]

	c.JSON(http.StatusOK, gin.H{
		"products": paginatedProducts,
		"query":    query,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// FilterProductsAPI filters products based on multiple criteria
func (papi *ProductAPI) FilterProductsAPI(c *gin.Context) {
	filters := bson.M{}

	// Build filters from query parameters
	if categoryID := c.Query("category_id"); categoryID != "" {
		catID, err := primitive.ObjectIDFromHex(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id format"})
			return
		}
		filters["$or"] = []bson.M{
			{"category_id": catID},
			{"categories": catID},
		}
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			if _, exists := filters["price"]; !exists {
				filters["price"] = bson.M{}
			}
			filters["price"].(bson.M)["$gte"] = min
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			if _, exists := filters["price"]; !exists {
				filters["price"] = bson.M{}
			}
			filters["price"].(bson.M)["$lte"] = max
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if brand := c.Query("brand"); brand != "" {
		filters["brand"] = primitive.Regex{Pattern: brand, Options: "i"}
	}

	if inStock := c.Query("in_stock"); inStock != "" {
		switch inStock {
		case "true":
			filters["stock_quantity"] = bson.M{"$gt": 0}
		case "false":
			filters["stock_quantity"] = bson.M{"$lte": 0}
		}
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := papi.svc.FilterProducts(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply pagination
	total := len(products)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedProducts := products[start:end]

	c.JSON(http.StatusOK, gin.H{
		"products": paginatedProducts,
		"filters":  filters,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}
