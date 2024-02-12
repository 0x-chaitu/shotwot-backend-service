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

type BriefsRepo struct {
	db *mongo.Collection
}

func NewBriefsRepo(db *mongo.Database) *BriefsRepo {
	return &BriefsRepo{
		db: db.Collection("brief"),
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
	cursor.SetBatchSize(20)

	var results []*domain.Brief
	if err = cursor.All(context.TODO(), &results); err != nil {
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
		return nil, err
	}
	updatedBrief := domain.Brief{}
	decodeErr := updatedResult.Decode(&updatedBrief)
	return &updatedBrief, decodeErr
}

func (r *BriefsRepo) Get(ctx context.Context, id string) (*domain.Brief, error) {
	filter := bson.M{"_id": id}
	result := r.db.FindOne(ctx, filter)
	if err := handleSingleError(result); err != nil {
		return nil, err
	}
	user := domain.User{}
	decodeErr := result.Decode(&user)
	return nil, decodeErr
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
