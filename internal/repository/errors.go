package repository

import (
	"shotwot_backend/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

func handleSingleError(result *mongo.SingleResult) error {
	if mongo.ErrNoDocuments == result.Err() {
		return domain.ErrNotFound
	} else if result.Err() != nil {
		return result.Err()
	}
	return nil
}
