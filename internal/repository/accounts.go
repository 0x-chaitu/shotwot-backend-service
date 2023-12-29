package repository

import (
	"context"
	"shotwot_backend/internal/domain"

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
	return nil, nil
}

func (r *AccountsRepo) GetByCredentials(ctx context.Context, userIdentifier, password string) (domain.Account, error) {
	return domain.Account{}, nil
}
