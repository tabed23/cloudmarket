package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/service/category"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryAPI struct {
	svc category.CategoryService
}

func NewCategoryAPI(svc category.CategoryService) *CategoryAPI {
	return &CategoryAPI{svc: svc}
}

func (capi *CategoryAPI) CreateCategoryAPI(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Call the service layer to create the category
	category, err := capi.svc.CreateCategory(c.Request.Context(), &category)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(201, category)

}

// GetCategoryAPI retrieves a single category by ID
func (capi *CategoryAPI) GetCategoryAPI(c *gin.Context) {
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id format"})
		return
	}

	category, err := capi.svc.GetCategoryByID(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if category exists (assuming empty category means not found)
	if category.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategoryAPI updates an existing category
func (capi *CategoryAPI) UpdateCategoryAPI(c *gin.Context) {
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id format"})
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent changing the ID
	category.ID = categoryID

	updatedCategory, err := capi.svc.UpdateCategory(c.Request.Context(), categoryID, &category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}
func (capi *CategoryAPI) DeleteCategoryAPI(c *gin.Context) {
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id format"})
		return
	}

	// Optional: Check if category has subcategories or products
	subcategories, err := capi.svc.GetSubcategories(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check subcategories"})
		return
	}

	if len(subcategories) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":               "cannot delete category with subcategories",
			"subcategories_count": len(subcategories),
		})
		return
	}

	err = capi.svc.DeleteCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

func (capi *CategoryAPI) ListCategoriesAPI(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Parse optional filters
	includeSubcategories := c.DefaultQuery("include_subcategories", "true") == "true"
	parentOnly := c.DefaultQuery("parent_only", "false") == "true"

	var categories []models.Category
	var err error

	if parentOnly {
		// Get only parent categories (categories without parent_id)
		filters := bson.M{
			"$or": []bson.M{
				{"parent_id": bson.M{"$exists": false}},
				{"parent_id": nil},
			},
		}
		categories, err = capi.svc.FilterCategories(c.Request.Context(), filters)
	} else {
		categories, err = capi.svc.GetAllCategories(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simple pagination logic
	total := len(categories)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedCategories := categories[start:end]

	// Optional: Add subcategory information
	if includeSubcategories {
		for i, category := range paginatedCategories {
			subcategories, err := capi.svc.GetSubcategories(c.Request.Context(), category.ID)
			if err == nil {
				// Add subcategory count or subcategories to response
				// This assumes you have a field for this in your model
				// You might need to create a response struct instead
				_ = subcategories // placeholder - you can modify based on your model
			}
			paginatedCategories[i] = category
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": paginatedCategories,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

func (capi *CategoryAPI) SearchCategoriesAPI(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	categories, err := capi.svc.SearchCategories(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// FilterCategoriesAPI filters categories based on query parameters
func (capi *CategoryAPI) FilterCategoriesAPI(c *gin.Context) {
	filters := bson.M{}

	// Build filters from query parameters
	if parentID := c.Query("parent_id"); parentID != "" {
		if parentID == "null" {
			filters["parent_id"] = bson.M{"$exists": false}
		} else {
			objID, err := primitive.ObjectIDFromHex(parentID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_id format"})
				return
			}
			filters["parent_id"] = objID
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if name := c.Query("name"); name != "" {
		filters["name"] = primitive.Regex{Pattern: name, Options: "i"}
	}

	categories, err := capi.svc.FilterCategories(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// AssignProductToCategoryAPI assigns a product to a category
func (capi *CategoryAPI) AssignProductToCategoryAPI(c *gin.Context) {
	type AssignRequest struct {
		ProductID  string `json:"product_id" binding:"required"`
		CategoryID string `json:"category_id" binding:"required"`
	}

	var req AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id format"})
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id format"})
		return
	}

	err = capi.svc.AssignProductToCategory(c.Request.Context(), productID, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product assigned to category successfully"})
}

// RemoveProductFromCategoryAPI removes a product from a category
func (capi *CategoryAPI) RemoveProductFromCategoryAPI(c *gin.Context) {
	type RemoveRequest struct {
		ProductID  string `json:"product_id" binding:"required"`
		CategoryID string `json:"category_id" binding:"required"`
	}

	var req RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id format"})
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id format"})
		return
	}

	err = capi.svc.RemoveProductFromCategory(c.Request.Context(), productID, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product removed from category successfully"})
}

// ListProductsInCategoryAPI lists all products in a specific category
func (capi *CategoryAPI) ListProductsInCategoryAPI(c *gin.Context) {
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id format"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := capi.svc.GetProductsInCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simple pagination logic
	start := (page - 1) * limit
	end := start + limit

	total := len(products)
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedProducts := products[start:end]

	c.JSON(http.StatusOK, gin.H{
		"products": paginatedProducts,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// ListCategoriesForProductAPI lists all categories for a specific product
func (capi *CategoryAPI) ListCategoriesForProductAPI(c *gin.Context) {
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

	categories, err := capi.svc.GetCategoriesForProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetCategoryHierarchyAPI gets the full hierarchy path for a category
func (capi *CategoryAPI) GetCategoryHierarchyAPI(c *gin.Context) {
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id format"})
		return
	}

	hierarchy, err := capi.svc.GetCategoryHierarchy(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hierarchy": hierarchy,
		"depth":     len(hierarchy),
	})
}
