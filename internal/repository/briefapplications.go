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

func (r *BriefApplicationsRepo) GetBriefApplications(ctx context.Context, predicate *helper.BriefApplicationsPredicate) ([]*domain.BriefApplication, error) {

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "isActive", Value: predicate.IsActive}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created", Value: predicate.Order}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user"},
			{Key: "localField", Value: "userId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "users"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "brief"},
			{Key: "localField", Value: "briefId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "briefs"},
		}}},
	}
	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.SetBatchSize(20)

	var results []*domain.BriefApplication
	for cursor.Next(ctx) {
		var briefApp struct {
			domain.BriefApplication
			Users  []domain.User  `bson:"users"`
			Briefs []domain.Brief `bson:"briefs"`
		}
		if err := cursor.Decode(&briefApp); err != nil {
			return nil, err
		}
		if len(briefApp.Users) > 0 {
			briefApp.User = briefApp.Users[0]
		}
		if len(briefApp.Briefs) > 0 {
			briefApp.Brief = briefApp.Briefs[0]
		}
		results = append(results, &briefApp.BriefApplication)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
