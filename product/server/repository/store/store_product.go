package store

import (
	"context"
	"fmt"
	"regexp"

	"github.com/tabed23/cloudmarket-product/server/models"
	"github.com/tabed23/cloudmarket-product/server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductStore struct {
	collection *mongo.Collection
}

func NewProductStore(collection *mongo.Collection) repository.ProductRepository {
	return &ProductStore{collection: collection}
}

// CreateProduct implements repository.ProductRepository.
func (p *ProductStore) CreateProduct(ctx context.Context, product *models.Product) (models.Product, error) {
	result, err := p.collection.InsertOne(ctx, product)
	if err != nil {
		return models.Product{}, err
	}
	product.ID = result.InsertedID.(primitive.ObjectID)
	return *product, nil
}

// DeleteProduct implements repository.ProductRepository.
func (p *ProductStore) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	_, err := p.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetAllProducts implements repository.ProductRepository.
func (p *ProductStore) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	cursor, err := p.collection.Find(ctx, primitive.M{})
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
		return nil, fmt.Errorf("cursor error: %v", err)
	}
	return products, nil
}

// GetProductByID implements repository.ProductRepository.
func (p *ProductStore) GetProductByID(ctx context.Context, id primitive.ObjectID) (models.Product, error) {
	var product models.Product
	err := p.collection.FindOne(ctx, primitive.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Product{}, nil
		}
		return models.Product{}, err
	}
	return product, nil
}

// GetProductsByCategory implements repository.ProductRepository.
func (p *ProductStore) GetProductsByCategory(ctx context.Context, categoryID primitive.ObjectID) ([]models.Product, error) {
	cursor, err := p.collection.Find(ctx, bson.M{"category_id": categoryID})
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

// UpdateProduct implements repository.ProductRepository.
func (p *ProductStore) UpdateProduct(ctx context.Context, id primitive.ObjectID, product *models.Product) (models.Product, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": product}
	_, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Product{}, err
	}
	return *product, nil
}

func (p *ProductStore) FilterProducts(ctx context.Context, filters bson.M) ([]models.Product, error) {
	cursor, err := p.collection.Find(ctx, filters)
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
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return products, nil
}

func (p *ProductStore) SearchProducts(ctx context.Context, query string) ([]models.Product, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": primitive.Regex{Pattern: regexp.QuoteMeta(query), Options: "i"}},
			{"description": primitive.Regex{Pattern: regexp.QuoteMeta(query), Options: "i"}},
			{"brand": primitive.Regex{Pattern: regexp.QuoteMeta(query), Options: "i"}},
			{"tags": primitive.Regex{Pattern: regexp.QuoteMeta(query), Options: "i"}},
		},
	}

	cursor, err := p.collection.Find(ctx, filter)
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
