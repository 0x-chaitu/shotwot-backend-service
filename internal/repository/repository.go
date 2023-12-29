package repository

import (
	"context"
	"shotwot_backend/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type Accounts interface {
	Create(ctx context.Context, user *domain.Account) (*domain.Account, error)

	//userIdentifier can be username, email
	GetByCredentials(ctx context.Context, userIdentifier, password string) (domain.Account, error)
}

type Repositories struct {
	Accounts Accounts
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Accounts: NewAccountsRepo(db),
	}
}
