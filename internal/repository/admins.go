package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/database/mongodb"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminsRepo struct {
	db *mongo.Collection
}

func NewAdminsRepo(db *mongo.Database) *AdminsRepo {
	return &AdminsRepo{
		db: db.Collection("admin"),
	}
}

func (r *AdminsRepo) Create(ctx context.Context, admin *domain.Admin) error {
	_, err := r.db.InsertOne(ctx, admin)
	if mongodb.IsDuplicate(err) {
		return domain.ErrAccountAlreadyExists
	}
	return nil
}

func (r *AdminsRepo) Update(ctx context.Context, admin *domain.Admin) (*domain.Admin, error) {
	result, err := helper.TODoc(admin)
	update := bson.M{
		"$set": result,
	}
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": admin.Id}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	updatedResult := r.db.FindOneAndUpdate(ctx, filter, update, &opt)
	if err := handleSingleError(updatedResult); err != nil {
		return nil, err
	}
	updatedadmin := domain.Admin{}
	decodeErr := updatedResult.Decode(&updatedadmin)
	return &updatedadmin, decodeErr
}

func (r *AdminsRepo) GetAdmins(ctx context.Context) ([]*domain.Admin, error) {

	cursor, err := r.db.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	// cursor.SetBatchSize(20)

	var results []*domain.Admin
	if err = cursor.All(context.TODO(), &results); err != nil {
		logger.Error(err)
		return nil, err
	}
	return results, nil
}

func (r *AdminsRepo) Get(ctx context.Context, id string) (*domain.Admin, error) {
	filter := bson.M{"_id": id}
	result := r.db.FindOne(ctx, filter)
	if err := handleSingleError(result); err != nil {
		logger.Error(err)
		return nil, err
	}
	admin := domain.Admin{}
	decodeErr := result.Decode(&admin)
	return &admin, decodeErr
}
