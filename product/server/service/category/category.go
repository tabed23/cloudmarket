package category

import (
	"context"
	"fmt"

	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/repository"
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
