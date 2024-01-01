package repository

import (
	"context"
	"errors"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/database/mongodb"
	"shotwot_backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection("Users"),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	_, err := r.db.InsertOne(ctx, user)
	logger.Info(user.Created)
	if mongodb.IsDuplicate(err) {
		return nil, domain.ErrAccountAlreadyExists
	}
	return nil, nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, userIdentifier, password string) (domain.User, error) {
	return domain.User{}, nil
}

func (r *UsersRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	update := bson.M{
		"$set": user,
	}
	filter := bson.M{"_id": user.Id}
	//get user details before update query
	result := r.db.FindOne(ctx, filter)
	if err := handleSingleError(result); err != nil {
		return nil, err
	}
	getUser := domain.User{}
	decodeErr := result.Decode(&getUser)
	if decodeErr != nil {
		return nil, decodeErr
	}

	if getUser.Email != user.Email {
		return nil, errors.New("user email invalid")
	} else if getUser.Pro != user.Pro {
		return nil, errors.New("user action invalid")
	} else if getUser.Created != user.Created {
		return nil, errors.New("user action invalid")
	} else if getUser.ProfileImage != user.ProfileImage {
		return nil, errors.New("user action invalid")
	}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	updatedResult := r.db.FindOneAndUpdate(ctx, filter, update, &opt)
	if err := handleSingleError(result); err != nil {
		return nil, err
	}
	updatedUser := domain.User{}
	decodeErr = updatedResult.Decode(&updatedUser)
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
