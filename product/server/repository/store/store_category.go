package store

import (
	"context"
	"regexp"

	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/repository"
	"github.com/tabed23/cloudmarket-product/server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryStore struct {
	categoryCollection *mongo.Collection
	productCollection  *mongo.Collection
}

func NewCategoryStore(categoryCollection, productCollection *mongo.Collection) repository.CategoryRepository {
	return &CategoryStore{categoryCollection: categoryCollection, productCollection: productCollection}
}

// CreateCategory implements repository.CategoryRepository.
func (c *CategoryStore) CreateCategory(ctx context.Context, category *models.Category) (models.Category, error) {
	result, err := c.categoryCollection.InsertOne(ctx, category)
	if err != nil {
		return models.Category{}, err
	}
	category.ID = result.InsertedID.(primitive.ObjectID)
	return *category, nil
}

// DeleteCategory implements repository.CategoryRepository.
func (c *CategoryStore) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	_, err := c.categoryCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetAllCategories implements repository.CategoryRepository.
func (c *CategoryStore) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	cursor, err := c.categoryCollection.Find(ctx, primitive.M{})
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
	err := c.categoryCollection.FindOne(ctx, primitive.M{"_id": id}).Decode(&category)
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
	cursor, err := c.categoryCollection.Find(ctx, bson.M{"parent_id": parentID})
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

	_, err := c.categoryCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Category{}, err
	}
	return *category, nil
}

func (c *CategoryStore) SearchCategories(ctx context.Context, query string) ([]models.Category, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": primitive.Regex{Pattern: regexp.QuoteMeta(query), Options: "i"}},
			{"description": primitive.Regex{Pattern: regexp.QuoteMeta(query), Options: "i"}},
		},
	}
	coursor, err := c.categoryCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer coursor.Close(ctx)
	var categories []models.Category
	for coursor.Next(ctx) {
		var category models.Category
		if err := coursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := coursor.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *CategoryStore) FilterCategories(ctx context.Context, filters bson.M) ([]models.Category, error) {
	courser, err := c.categoryCollection.Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	defer courser.Close(ctx)
	var categories []models.Category
	for courser.Next(ctx) {
		var category models.Category
		if err := courser.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := courser.Err(); err != nil {
		return nil, err
	}
	return categories, nil

}

func (c *CategoryStore) AssignProductToCategory(ctx context.Context, productID, categoryID primitive.ObjectID) error {
	filter := bson.M{"_id": productID}
	update := bson.M{
		"$addToSet": bson.M{
			"category_ids": categoryID},
	}
	_, err := c.productCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	categoryFilter := bson.M{"_id": categoryID}
	categoryUpdate := bson.M{
		"$addToSet": bson.M{"product_ids": productID}}
	_, err = c.categoryCollection.UpdateOne(ctx, categoryFilter, categoryUpdate)
	return err
}
func (c *CategoryStore) RemoveProductFromCategory(ctx context.Context, productID, categoryID primitive.ObjectID) error {
	filter := bson.M{"_id": productID}
	update := bson.M{
		"$pull": bson.M{
			"category_ids": categoryID},
	}
	_, err := c.productCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	categoryFilter := bson.M{"_id": categoryID}
	categoryUpdate := bson.M{
		"$pull": bson.M{"product_ids": productID}}
	_, err = c.categoryCollection.UpdateOne(ctx, categoryFilter, categoryUpdate)
	return err
}
func (c *CategoryStore) GetProductsInCategory(ctx context.Context, categoryID primitive.ObjectID) ([]models.Product, error) {
	filter := bson.M{"categories": categoryID}
	cursor, err := c.productCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var products []models.Product
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (c *CategoryStore) GetCategoryHierarchy(ctx context.Context, categoryID primitive.ObjectID) ([]models.Category, error) {
	var hierarchy []models.Category
	currentID := categoryID

	for !currentID.IsZero() {
		var category models.Category
		err := c.categoryCollection.FindOne(ctx, bson.M{"_id": currentID}).Decode(&category)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				break
			}
			return nil, err
		}
		hierarchy = append([]models.Category{category}, hierarchy...)
		if category.ParentID != nil && !category.ParentID.IsZero() {
			currentID = *category.ParentID
		} else {
			break
		}

	}
	return hierarchy, nil
}

func (c *CategoryStore) GetCategoriesForProduct(ctx context.Context, productID primitive.ObjectID) ([]models.Category, error) {
	// Get the product
	var product models.Product
	err := c.productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Category{}, err
		}
		return nil, err
	}

	// Collect all category IDs (primary + additional categories)
	var categoryIDs []primitive.ObjectID

	// Add primary category if exists
	if !product.CategoryID.IsZero() {
		categoryIDs = append(categoryIDs, product.CategoryID)
	}

	// Add additional categories if any
	categoryIDs = append(categoryIDs, product.Categories...)

	// Remove duplicates
	categoryIDs = utils.RemoveDuplicateObjectIDs(categoryIDs)

	if len(categoryIDs) == 0 {
		return []models.Category{}, nil
	}

	// Find all categories
	cursor, err := c.categoryCollection.Find(ctx, bson.M{"_id": bson.M{"$in": categoryIDs}})
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
