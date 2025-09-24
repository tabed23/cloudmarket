package repository

import (
	"context"

	"github.com/tabed23/cloudmarket-product/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) (models.Product, error)
	GetProductByID(ctx context.Context, id primitive.ObjectID) (models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	UpdateProduct(ctx context.Context, id primitive.ObjectID, product *models.Product) (models.Product, error)
	DeleteProduct(ctx context.Context, id primitive.ObjectID) error
	GetProductsByCategory(ctx context.Context, categoryID primitive.ObjectID) ([]models.Product, error)
	FilterProducts(ctx context.Context, filters bson.M) ([]models.Product, error)

}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *models.Category) (models.Category, error)
	GetCategoryByID(ctx context.Context, id primitive.ObjectID) (models.Category, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	UpdateCategory(ctx context.Context, id primitive.ObjectID, category *models.Category) (models.Category, error)
	DeleteCategory(ctx context.Context, id primitive.ObjectID) error
	GetSubcategories(ctx context.Context, parentID primitive.ObjectID) ([]models.Category, error)
}