package product

import (
	"context"
	"fmt"

	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	repo repository.ProductRepository
}

// CreateProduct implements repository.ProductRepository.
func (p *ProductService) CreateProduct(ctx context.Context, product *models.Product) (models.Product, error) {
	createProduct, err := p.repo.CreateProduct(ctx, product)
	if err != nil {
		return models.Product{}, err
	}
	return createProduct, nil
}

// DeleteProduct implements repository.ProductRepository.
func (p *ProductService) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	err := p.repo.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting product: %v", err)
	}

	return nil
}

// FilterProducts implements repository.ProductRepository.
func (p *ProductService) FilterProducts(ctx context.Context, filters bson.M) ([]models.Product, error) {
	bsonFilters := make(map[string]interface{})

	// Convert map to BSON filters
	for key, value := range filters {
		bsonFilters[key] = value
	}

	products, err := p.repo.FilterProducts(ctx, bsonFilters)
	if err != nil {
		return nil, fmt.Errorf("error applying filters: %v", err)
	}

	return products, nil
}

// GetAllProducts implements repository.ProductRepository.
func (p *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := p.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %v", err)
	}

	return products, nil
}

// GetProductByID implements repository.ProductRepository.
func (p *ProductService) GetProductByID(ctx context.Context, id primitive.ObjectID) (models.Product, error) {
	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		return models.Product{}, fmt.Errorf("error fetching product: %v", err)
	}

	return product, nil
}

// GetProductsByCategory implements repository.ProductRepository.
func (p *ProductService) GetProductsByCategory(ctx context.Context, categoryID primitive.ObjectID) ([]models.Product, error) {

	products, err := p.repo.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error fetching products by category: %v", err)
	}

	return products, nil
}

// UpdateProduct implements repository.ProductRepository.
func (p *ProductService) UpdateProduct(ctx context.Context, id primitive.ObjectID, product *models.Product) (models.Product, error) {
	updatedProduct, err := p.repo.UpdateProduct(ctx, id, product)
	if err != nil {
		return models.Product{}, fmt.Errorf("error updating product: %v", err)
	}

	return updatedProduct, nil
}
func (p *ProductService) SearchProducts(ctx context.Context, query string) ([]models.Product, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	products, err := p.repo.SearchProducts(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error searching products: %v", err)
	}
	return products, nil
}

func NewProductService(repo repository.ProductRepository) repository.ProductRepository {
	return &ProductService{repo: repo}
}
