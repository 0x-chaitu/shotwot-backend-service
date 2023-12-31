package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/database/mongodb"
	"shotwot_backend/pkg/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection("user"),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return domain.ErrAccountAlreadyExists
	}
	return nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, userIdentifier, password string) (domain.User, error) {
	return domain.User{}, nil
}

func (r *UsersRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	filter := bson.M{"_id": user.Id}

	result, err := helper.TODoc(user)
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
	updatedUser := domain.User{}
	decodeErr := updatedResult.Decode(&updatedUser)
	return &updatedUser, decodeErr
}

func (r *UsersRepo) Get(ctx context.Context, id string) (*domain.User, error) {
	filter := bson.M{"_id": id}
	result := r.db.FindOne(ctx, filter)
	if err := handleSingleError(result); err != nil {
		return nil, err
	}
	user := domain.User{}
	decodeErr := result.Decode(&user)
	return &user, decodeErr
}
