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

func (r *BriefApplicationsRepo) Update(ctx context.Context, application *domain.BriefApplication) (*domain.BriefApplication, error) {
	filter := bson.M{"_id": application.Id}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: application.Status}}}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	updatedResult := r.db.FindOneAndUpdate(ctx, filter, update, &opt)
	if err := handleSingleError(updatedResult); err != nil {
		return nil, err
	}
	updatedBrief := domain.BriefApplication{}
	decodeErr := updatedResult.Decode(&updatedBrief)
	return &updatedBrief, decodeErr
}

func (r *BriefApplicationsRepo) GetBriefApplications(ctx context.Context, id string) ([]*domain.BriefApplication, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "briefId",
			Value: objID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user"},
			{Key: "localField", Value: "userId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: "$user"}},
		{{
			Key: "$project", Value: bson.M{
				"user": bson.M{
					"username": 1,
				},
				"briefId": 1,
				"created": 1,
			},
		}},
	}
	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.SetBatchSize(20)

	var results []*domain.BriefApplication
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *BriefApplicationsRepo) GetBriefApplication(ctx context.Context, id string) (*domain.UserBriefAppliedDetails, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id",
			Value: objID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user"},
			{Key: "localField", Value: "userId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: "$user"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "brief"},
			{Key: "localField", Value: "briefId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "brief"},
		}}},
		{{Key: "$unwind", Value: "$brief"}},
	}
	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.SetBatchSize(20)

	var result []*domain.UserBriefAppliedDetails
	if err = cursor.All(ctx, &result); err != nil {
		logger.Error(err)
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil

}

func (r *BriefApplicationsRepo) GetUserBriefApplications(ctx context.Context, id string) ([]*domain.UserBriefAppliedDetails, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "userId",
			Value: id}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "brief"},
			{Key: "localField", Value: "briefId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "brief"},
		}}},
		{{Key: "$unwind", Value: "$brief"}},
		{{Key: "$sort", Value: bson.M{"created": -1}}},
	}
	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.SetBatchSize(20)

	var result []*domain.UserBriefAppliedDetails
	if err = cursor.All(ctx, &result); err != nil {
		logger.Error(err)
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	return nil, nil

}
