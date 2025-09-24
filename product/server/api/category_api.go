package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/service/category"
)

type CategoryAPI struct {
	svc category.CategoryService
}
func NewCategoryAPI(svc category.CategoryService) *CategoryAPI {
	return &CategoryAPI{svc: svc}
}

func (capi *CategoryAPI)CreateCategoryAPI(c *gin.Context) {
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
func (capi *CategoryAPI)GetCategoryAPI(c *gin.Context) {}
func (capi *CategoryAPI)UpdateCategoryAPI(c *gin.Context) {}
func (capi *CategoryAPI)DeleteCategoryAPI(c *gin.Context) {}
func (capi *CategoryAPI)ListCategoriesAPI(c *gin.Context) {}
func (capi *CategoryAPI)SearchCategoriesAPI(c *gin.Context) {}
func (capi *CategoryAPI)FilterCategoriesAPI(c *gin.Context) {}
func (capi *CategoryAPI)AssignProductToCategoryAPI(c *gin.Context) {}
func (capi *CategoryAPI)RemoveProductFromCategoryAPI(c *gin.Context) {}
func (capi *CategoryAPI)ListProductsInCategoryAPI(c *gin.Context) {}
func (capi *CategoryAPI)ListCategoriesForProductAPI(c *gin.Context) {}
func (capi *CategoryAPI)GetCategoryHierarchyAPI(c *gin.Context) {}
