package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/helper"
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

func (r *BriefApplicationsRepo) GetBriefApplications(ctx context.Context, predicate *helper.BriefApplicationsPredicate) ([]*domain.BriefApplication, error) {
	filter := bson.D{}
	if predicate.IsActive != nil {
		logger.Info(predicate.IsActive)
		filter = append(filter, bson.E{Key: "isActive", Value: predicate.IsActive})

	}
	opts := options.Find().SetSort(bson.D{{Key: "created", Value: predicate.Order}})
	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.SetBatchSize(20)

	var results []*domain.BriefApplication
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
