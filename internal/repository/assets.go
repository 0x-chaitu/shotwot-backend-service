package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssetsRepo struct {
	db *mongo.Collection
}

func NewAssetRepo(db *mongo.Database) *AssetsRepo {
	return &AssetsRepo{
		db: db.Collection("asset"),
	}
}

// Create Brief Application Api
func (r *AssetsRepo) Create(ctx context.Context, asset *domain.Asset) (*domain.Asset, error) {
	result, err := r.db.InsertOne(ctx, asset)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	asset.Id = result.InsertedID.(primitive.ObjectID)
	return asset, nil
}
