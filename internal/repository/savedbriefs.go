package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SavedBriefsRepo struct {
	db *mongo.Collection
}

func NewSavedBriefsRepo(db *mongo.Database) *SavedBriefsRepo {
	return &SavedBriefsRepo{
		db: db.Collection("savedbriefs"),
	}
}

// Create and Update
func (r *SavedBriefsRepo) CreateOrUpdate(ctx context.Context, saveBrief *domain.SavedBrief) (*domain.SavedBrief, error) {
	// Check if a saved brief with the given userId and briefId already exists
	filter := bson.M{"userId": saveBrief.UserId, "briefId": saveBrief.BriefId}
	existingBrief := &domain.SavedBrief{}
	err := r.db.FindOne(ctx, filter).Decode(existingBrief)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No existing saved brief found, set status to true and insert the new one
			saveBrief.Status = true
			return r.Create(ctx, saveBrief)
		}
		// Some other error occurred
		logger.Error(err)
		return nil, err
	}

	// Existing saved brief found, toggle the status
	saveBrief.Status = !existingBrief.Status

	// Update the status of the existing saved brief
	update := bson.M{"$set": bson.M{"status": saveBrief.Status}}
	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// Return the updated saved brief
	return existingBrief, nil
}
func (r *SavedBriefsRepo) Create(ctx context.Context, saveBrief *domain.SavedBrief) (*domain.SavedBrief, error) {
	result, err := r.db.InsertOne(ctx, saveBrief)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	saveBrief.Id = result.InsertedID.(primitive.ObjectID)
	return saveBrief, nil
}

// Get All Saved Briefs
// GetSavedBriefs retrieves saved briefs by user ID
func (r *SavedBriefsRepo) GetSavedBriefs(ctx context.Context, userId string) ([]*domain.SavedBrief, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "userId", Value: userId},
		{Key: "status", Value: true},
	}}}

	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "brief"},
		{Key: "localField", Value: "briefId"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "brief"},
	}}}

	// Unwind the brief array to deconstruct the array elements
	unwindStage := bson.D{{Key: "$unwind", Value: "$brief"}}

	// Project the necessary fields to reshape the output
	projectStage := bson.D{{Key: "$project", Value: bson.M{
		"_id":     1,
		"briefId": 1,
		"userId":  1,
		"created": 1,
		"updated": 1,
		"status":  1,
		"brief":   "$brief",
	}}}

	// Aggregation pipeline
	pipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage, projectStage}

	// Execute aggregation
	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var savedBriefs []*domain.SavedBrief
	for cursor.Next(ctx) {
		var savedBrief domain.SavedBrief
		if err := cursor.Decode(&savedBrief); err != nil {
			return nil, err
		}
		savedBriefs = append(savedBriefs, &savedBrief)
	}

	return savedBriefs, nil
}
