package store

import (
	"context"

	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryStore struct {
	collection *mongo.Collection
}

func NewCategoryStore(collection *mongo.Collection) repository.CategoryRepository {
	return &CategoryStore{collection: collection}
}

// CreateCategory implements repository.CategoryRepository.
func (c *CategoryStore) CreateCategory(ctx context.Context, category *models.Category) (models.Category, error) {
	result,err := c.collection.InsertOne(ctx,category)
	if err != nil {
		return models.Category{}, err
	}
	category.ID = result.InsertedID.(primitive.ObjectID)
	return *category, nil
}

// DeleteCategory implements repository.CategoryRepository.
func (c *CategoryStore) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	_, err := c.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetAllCategories implements repository.CategoryRepository.
func (c *CategoryStore) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	cursor, err := c.collection.Find(ctx, primitive.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	for cursor.Next(ctx) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID implements repository.CategoryRepository.
func (c *CategoryStore) GetCategoryByID(ctx context.Context, id primitive.ObjectID) (models.Category, error) {
	var category models.Category
	err := c.collection.FindOne(ctx, primitive.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Category{}, nil
		}
		return models.Category{}, err
	}
	return category, nil
}

// GetSubcategories implements repository.CategoryRepository.
func (c *CategoryStore) GetSubcategories(ctx context.Context, parentID primitive.ObjectID) ([]models.Category, error) {
	cursor, err := c.collection.Find(ctx, bson.M{"parent_id": parentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	for cursor.Next(ctx) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// UpdateCategory implements repository.CategoryRepository.
func (c *CategoryStore) UpdateCategory(ctx context.Context, id primitive.ObjectID, category *models.Category) (models.Category, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": category}

	_, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Category{}, err
	}
	return *category, nil
}

