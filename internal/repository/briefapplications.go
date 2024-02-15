package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BriefApplicationsRepo struct {
	db *mongo.Collection
}

func NewBriefApplicationsRepo(db *mongo.Database) *BriefApplicationsRepo {
	return &BriefApplicationsRepo{
		db: db.Collection("briefapplication"),
	}
}

// Create Brief Application Api
func (r *BriefApplicationsRepo) Create(ctx context.Context, briefapplication *domain.BriefApplication) (*domain.BriefApplication, error) {
	result, err := r.db.InsertOne(ctx, briefapplication)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	briefapplication.Id = result.InsertedID.(primitive.ObjectID)
	return briefapplication, nil
}

func (r *BriefApplicationsRepo) GetBriefApplications(ctx context.Context, id string) ([]*domain.BriefApplication, error) {

	// pipeline := mongo.Pipeline{
	// bson.D{{Key: "$lookup", Value: bson.D{
	// 	{Key: "from", Value: "brief"},
	// 	{Key: "localField", Value: "briefId"},
	// 	{Key: "foreignField", Value: "_id"},
	// 	{Key: "as", Value: "briefs"},
	// }}},
	// bson.D{{Key: "$match", Value: bson.D{{Key: "isActive", Value: predicate.IsActive}}}},
	// bson.D{{Key: "$sort", Value: bson.D{{Key: "created", Value: predicate.Order}}}},
	// bson.D{{Key: "$lookup", Value: bson.D{
	// 	{Key: "from", Value: "user"},
	// 	{Key: "localField", Value: "userId"},
	// 	{Key: "foreignField", Value: "_id"},
	// 	{Key: "as", Value: "users"},
	// }}},
	// }
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{
		Key:   "briefId",
		Value: objID,
	}}
	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})
	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.SetBatchSize(20)

	var results []*domain.BriefApplication
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	logger.Info(results)
	return results, nil
}
