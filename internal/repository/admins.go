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

func (r *AdminsRepo) Get(ctx context.Context, id string) (*domain.Admin, error) {
	filter := bson.M{"_id": "Ze0EvtlUPmb8H3tYtkuK1sTDpgU2"}
	result := r.db.FindOne(ctx, filter)
	if err := handleSingleError(result); err != nil {
		logger.Error(err)
		return nil, err
	}
	admin := domain.Admin{}
	decodeErr := result.Decode(&admin)
	return &admin, decodeErr
}
