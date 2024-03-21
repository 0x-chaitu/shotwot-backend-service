package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
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

func (r *AssetsRepo) Update(ctx context.Context, asset *domain.Asset) (*domain.Asset, error) {
	result, err := helper.TODoc(asset)
	update := bson.M{
		"$set": result,
	}
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": asset.Id}

	res := r.db.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		logger.Error(err)
		return nil, res.Err()
	}

	return asset, nil
}

func (r *AssetsRepo) GetAssets(ctx context.Context) ([]*domain.Asset, error) {
	// var cond = "$lt"
	filter := primitive.D{}
	// if !predicate.StartDate.IsZero() {
	// 	filter = append(filter, bson.E{Key: "created", Value: bson.D{
	// 		{Key: cond, Value: predicate.StartDate}}})
	// } else {
	// 	filter = append(filter, bson.E{Key: "created", Value: bson.D{
	// 		{Key: cond, Value: time.Now()}}})
	// }

	// opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})
	// if predicate.Skip != 0 {
	// 	opts.SetSkip(int64(predicate.Skip) * 20)
	// }

	// opts.SetLimit(int64(20))
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	// cursor.SetBatchSize(100)

	results := []*domain.Asset{}

	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
