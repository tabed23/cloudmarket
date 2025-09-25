package category

import (
	"context"
	"fmt"

	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryService defines the methods for category-related operations
type CategoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(repo repository.CategoryRepository) repository.CategoryRepository {
	return &CategoryService{
		repo: repo,
	}
}


// CreateCategory implements repository.CategoryRepository.
func (c *CategoryService) CreateCategory(ctx context.Context, category *models.Category) (models.Category, error) {
	createCategory, err := c.repo.CreateCategory(ctx, category)
	if err != nil {
		return models.Category{}, fmt.Errorf("error creating category: %v", err)
	}
	return  createCategory, nil
}

// DeleteCategory implements repository.CategoryRepository.
func (c *CategoryService) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	err := c.repo.DeleteCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting category: %v", err)
	}
	return nil
}

// GetAllCategories implements repository.CategoryRepository.
func (c *CategoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := c.repo.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %v", err)
	}
	return categories, nil
}

// GetCategoryByID implements repository.CategoryRepository.
func (c *CategoryService) GetCategoryByID(ctx context.Context, id primitive.ObjectID) (models.Category, error) {
	category, err := c.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return models.Category{}, fmt.Errorf("error fetching category by ID: %v", err)
	}
	return category, nil
}

// GetSubcategories implements repository.CategoryRepository.
func (c *CategoryService) GetSubcategories(ctx context.Context, parentID primitive.ObjectID) ([]models.Category, error) {
	subcategories, err := c.repo.GetSubcategories(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("error fetching subcategories: %v", err)
	}
	return subcategories, nil
}

// UpdateCategory implements repository.CategoryRepository.
func (c *CategoryService) UpdateCategory(ctx context.Context, id primitive.ObjectID, category *models.Category) (models.Category, error) {
	updatedCategory, err := c.repo.UpdateCategory(ctx, id, category)
	if err != nil {
		return models.Category{}, fmt.Errorf("error updating category: %v", err)
	}
	return updatedCategory, nil
}

func (c *CategoryService) SearchCategories(ctx context.Context, query string) ([]models.Category, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	categories, err := c.repo.SearchCategories(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error searching categories: %v", err)
	}
	return categories, nil
}

// FilterCategories filters categories based on provided criteria
func (c *CategoryService) FilterCategories(ctx context.Context, filters bson.M) ([]models.Category, error) {
	categories, err := c.repo.FilterCategories(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error filtering categories: %v", err)
	}
	return categories, nil
}

// AssignProductToCategory assigns a product to a category
func (c *CategoryService) AssignProductToCategory(ctx context.Context, productID, categoryID primitive.ObjectID) error {
	err := c.repo.AssignProductToCategory(ctx, productID, categoryID)
	if err != nil {
		return fmt.Errorf("error assigning product to category: %v", err)
	}
	return nil
}

// RemoveProductFromCategory removes a product from a category
func (c *CategoryService) RemoveProductFromCategory(ctx context.Context, productID, categoryID primitive.ObjectID) error {
	err := c.repo.RemoveProductFromCategory(ctx, productID, categoryID)
	if err != nil {
		return fmt.Errorf("error removing product from category: %v", err)
	}
	return nil
}

// GetProductsInCategory gets all products in a specific category
func (c *CategoryService) GetProductsInCategory(ctx context.Context, categoryID primitive.ObjectID) ([]models.Product, error) {
	products, err := c.repo.GetProductsInCategory(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error fetching products in category: %v", err)
	}
	return products, nil
}

// GetCategoriesForProduct gets all categories for a specific product
func (c *CategoryService) GetCategoriesForProduct(ctx context.Context, productID primitive.ObjectID) ([]models.Category, error) {
	categories, err := c.repo.GetCategoriesForProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories for product: %v", err)
	}
	return categories, nil
}

// GetCategoryHierarchy gets the full hierarchy path for a category
func (c *CategoryService) GetCategoryHierarchy(ctx context.Context, categoryID primitive.ObjectID) ([]models.Category, error) {
	hierarchy, err := c.repo.GetCategoryHierarchy(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error fetching category hierarchy: %v", err)
	}
	return hierarchy, nil
}
