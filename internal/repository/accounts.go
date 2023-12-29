package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/database/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type AccountsRepo struct {
	db *mongo.Collection
}

func NewAccountsRepo(db *mongo.Database) *AccountsRepo {
	return &AccountsRepo{
		db: db.Collection("accounts"),
	}
}

func (r *AccountsRepo) Create(ctx context.Context, user *domain.Account) (*domain.Account, error) {
	_, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return nil, domain.ErrAccountAlreadyExists
	}
	return nil, nil
}

func (r *AccountsRepo) GetByCredentials(ctx context.Context, userIdentifier, password string) (domain.Account, error) {
	return domain.Account{}, nil
}
