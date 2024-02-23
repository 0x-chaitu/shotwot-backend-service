package repository

import (
	"context"
	"fmt"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MasterClassRepo struct {
	db *mongo.Collection
}

func NewMasterClassRepo(db *mongo.Database) *MasterClassRepo {
	return &MasterClassRepo{
		db: db.Collection("playlist"),
	}
}

func (r *MasterClassRepo) CreatePlaylist(ctx context.Context, playlist *domain.Playlist) (*domain.Playlist, error) {
	result, err := r.db.InsertOne(ctx, playlist)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	playlist.Id = insertedID.Hex() // Convert ObjectID to string

	return playlist, nil
}

func (r *MasterClassRepo) GetPlaylists(ctx context.Context, predicate *helper.PlaylistPredicate) ([]*domain.Playlist, error) {
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

	var results []*domain.Playlist
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
