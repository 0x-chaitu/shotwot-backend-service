package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BriefsRepo struct {
	db          *mongo.Collection
	archievedDb *mongo.Collection
}

func NewBriefsRepo(db *mongo.Database) *BriefsRepo {
	return &BriefsRepo{
		db:          db.Collection("brief"),
		archievedDb: db.Collection("draftBrief"),
	}
}

func (r *BriefsRepo) Create(ctx context.Context, brief *domain.Brief) (*domain.Brief, error) {
	result, err := r.db.InsertOne(ctx, brief)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	brief.Id = result.InsertedID.(primitive.ObjectID)
	return brief, nil
}

func (r *BriefsRepo) CreateDraft(ctx context.Context, brief *domain.Brief) (*domain.Brief, error) {
	result, err := r.archievedDb.InsertOne(ctx, brief)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	brief.Id = result.InsertedID.(primitive.ObjectID)
	return brief, nil
}

func (r *BriefsRepo) GetBriefs(ctx context.Context, predicate *helper.BriefPredicate) ([]*domain.Brief, error) {
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
	// cursor.SetBatchSize(20)

	var results []*domain.Brief
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BriefsRepo) GetUserBriefs(ctx context.Context, predicate *helper.BriefPredicate) ([]*domain.Brief, error) {
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "isActive", Value: true})

	if len(predicate.Type) > 0 {
		typeF := []bson.M{}
		for _, v := range predicate.Type {
			typeF = append(typeF, bson.M{"type": v})
		}
		filter = append(filter, bson.E{
			Key: "$or", Value: typeF},
		)
	}

	switch predicate.Expiry {
	case 1:
		filter = append(filter, bson.E{Key: "expiry", Value: bson.D{
			{Key: "$gt", Value: time.Now()}}})
		filter = append(filter, bson.E{Key: "expiry", Value: bson.D{
			{Key: "$lt", Value: time.Now().Add(24 * 7 * time.Hour)}}})
	case 2:
		filter = append(filter, bson.E{Key: "expiry", Value: bson.D{
			{Key: "$gt", Value: time.Now()}}})
		filter = append(filter, bson.E{Key: "expiry", Value: bson.D{
			{Key: "$lt", Value: time.Now().Add(24 * 15 * time.Hour)}}})
	case 3:
		filter = append(filter, bson.E{Key: "expiry", Value: bson.D{
			{Key: "$gt", Value: time.Now().Add(24 * 15 * time.Hour)}}})
		filter = append(filter, bson.E{Key: "expiry", Value: bson.D{
			{Key: "$lt", Value: time.Now().Add(24 * 30 * time.Hour)}}})
	case 4:
		filter = append(filter, bson.E{Key: "expiry", Value: bson.D{
			{Key: "$gt", Value: time.Now().Add(24 * 30 * time.Hour)}}})

	default:
	}

	if predicate.RewardG != 0 {
		filter = append(filter, bson.E{Key: "reward", Value: bson.D{
			{Key: "$gt", Value: predicate.RewardG}}})
	}

	if predicate.RewardL != 0 {
		filter = append(filter, bson.E{Key: "reward", Value: bson.D{
			{Key: "$lt", Value: predicate.RewardL}}})
	}

	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})
	// .SetSkip(predicate.Skip).SetLimit(5)
	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	// cursor.SetBatchSize(20)

	var results []*domain.Brief
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BriefsRepo) Update(ctx context.Context, brief *domain.Brief) (*domain.Brief, error) {
	filter := bson.M{"_id": brief.Id}

	result, err := helper.TODoc(brief)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": result,
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	updatedResult := r.db.FindOneAndUpdate(ctx, filter, update, &opt)
	if err := handleSingleError(updatedResult); err != nil {
		logger.Info(err)
		return nil, err
	}
	updatedBrief := domain.Brief{}
	decodeErr := updatedResult.Decode(&updatedBrief)
	return &updatedBrief, decodeErr
}

func (r *BriefsRepo) GetBrief(ctx context.Context, id string) (*domain.Brief, error) {
	briedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": briedId}
	result := r.db.FindOne(ctx, filter)

	brief := domain.Brief{}
	decodeErr := result.Decode(&brief)
	return &brief, decodeErr
}

func (r *BriefsRepo) DeleteBrief(ctx context.Context, id string) error {
	briedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": briedId}
	res := r.db.FindOneAndDelete(ctx, filter)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
