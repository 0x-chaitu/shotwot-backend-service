package repository

import (
	"context"
	"fmt"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/database/mongodb"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *UsersRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	res, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return nil, domain.ErrAccountAlreadyExists
	}
	if err != nil {
		return nil, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert")
	}
	user.Id = insertedID
	return user, nil
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
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": userId}
	result := r.db.FindOne(ctx, filter)
	if err := handleSingleError(result); err != nil {
		return nil, err
	}
	user := domain.User{}
	decodeErr := result.Decode(&user)
	return &user, decodeErr
}

func (r *UsersRepo) GetOrCreate(ctx context.Context, user *domain.User) (*domain.User, error) {
	filter := bson.M{"userId": user.UserId}
	result := r.db.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			res, err := r.db.InsertOne(ctx, user)
			if err != nil {
				return nil, err
			}
			insertedID, ok := res.InsertedID.(primitive.ObjectID)
			if !ok {
				return nil, fmt.Errorf("failed to convert")
			}
			user.Id = insertedID
			return user, err
		} else {
			return nil, result.Err()
		}
	}
	decodeErr := result.Decode(user)
	return user, decodeErr
}

func (r *UsersRepo) GetOrCreateByPhone(ctx context.Context, user *domain.User) (*domain.User, error) {
	filter := bson.M{"mobile": user.Mobile}
	result := r.db.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			res, err := r.db.InsertOne(ctx, user)
			if err != nil {
				return nil, err
			}
			insertedID, ok := res.InsertedID.(primitive.ObjectID)
			if !ok {
				return nil, fmt.Errorf("failed to convert")
			}
			user.Id = insertedID
			return user, err
		} else {
			return nil, result.Err()
		}
	}
	decodeErr := result.Decode(user)
	return user, decodeErr
}

func (r *UsersRepo) GetOrCreateByEmail(ctx context.Context, user *domain.User) (*domain.User, error) {
	filter := bson.M{"email": user.Email}
	result := r.db.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			res, err := r.db.InsertOne(ctx, user)
			if err != nil {
				return nil, err
			}
			insertedID, ok := res.InsertedID.(primitive.ObjectID)
			if !ok {
				return nil, fmt.Errorf("failed to convert")
			}
			user.Id = insertedID
			return user, err
		} else {
			return nil, result.Err()
		}
	}
	decodeErr := result.Decode(user)
	return user, decodeErr
}

func (r *UsersRepo) Download(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})
	// filter := bson.M{"created": bson.M{"$gte": predicate.StartDate,
	// 	"$lte": predicate.EndDate}}
	cursor, err := r.db.Find(ctx, primitive.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := []*domain.User{}

	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *UsersRepo) SearchUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error) {
	searchStage := bson.D{
		{Key: "$search", Value: bson.M{
			"index": "SearchUsers",
			"compound": bson.D{
				{Key: "should", Value: bson.A{
					bson.D{
						{Key: "autocomplete", Value: bson.D{{Key: "query", Value: predicate.Key}, {Key: "path", Value: "email"}}},
					},
					bson.D{
						{Key: "phrase", Value: bson.M{
							"query": predicate.Key,
							"path":  "_id"},
						}},
				},
				},
			}},
		}}

	limitStage := bson.D{{Key: "$limit", Value: 20}}
	skipStage := bson.D{{Key: "$skip", Value: predicate.Skip}}

	opts := options.Aggregate().SetMaxTime(5 * time.Second)

	cursor, err := r.db.Aggregate(ctx,
		mongo.Pipeline{searchStage, skipStage, limitStage}, opts)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var results []*domain.User

	if err = cursor.All(ctx, &results); err != nil {
		logger.Error(err)
		return nil, err
	}
	return results, err
}

func (r *UsersRepo) GetUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error) {
	var cond = "$lt"
	var filter primitive.D
	filter = primitive.D{}
	if !predicate.StartDate.IsZero() {
		filter = append(filter, bson.E{Key: "created", Value: bson.D{
			{Key: cond, Value: predicate.StartDate}}})
	} else {
		filter = append(filter, bson.E{Key: "created", Value: bson.D{
			{Key: cond, Value: time.Now()}}})
	}

	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})
	if predicate.Skip != 0 {
		opts.SetSkip(int64(predicate.Skip) * 20)
	}

	opts.SetLimit(int64(20))
	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	// cursor.SetBatchSize(100)

	results := []*domain.User{}

	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *UsersRepo) TotalUsers(ctx context.Context) (int64, error) {
	opts := options.Count().SetHint("_id_")
	count, err := r.db.CountDocuments(ctx, bson.D{}, opts)
	if err != nil {
		return 0, err
	}
	return count, nil
}
