package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tabed23/cloudmarket-product/server/service/product"
)

type ProductAPI struct {
	svc product.ProductService
}
func NewProductAPI(svc product.ProductService) *ProductAPI {
	return &ProductAPI{svc: svc}
}
func (papi *ProductAPI)CreateProductAPI(c *gin.Context) {
	// Placeholder for creating product API endpoints
}
func (papi *ProductAPI)GetProductAPI(c *gin.Context) {}
func (papi *ProductAPI)UpdateProductAPI(c *gin.Context) {}
func (papi *ProductAPI)DeleteProductAPI(c *gin.Context) {}
func (papi *ProductAPI)ListProductsAPI(c *gin.Context) {}
func (papi *ProductAPI)SearchProductsAPI(c *gin.Context) {}
func (papi *ProductAPI)FilterProductsAPI(c *gin.Context) {}