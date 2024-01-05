package repository

import (
	"context"
	"shotwot_backend/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user *domain.User) error

	Update(ctx context.Context, user *domain.User) (*domain.User, error)

	Get(ctx context.Context, id string) (*domain.User, error)
}

type Admins interface {
	Create(ctx context.Context, admin *domain.Admin) (*domain.Admin, error)

	Update(ctx context.Context, admin *domain.Admin) (*domain.Admin, error)

	Get(ctx context.Context, id string) (*domain.Admin, error)
}

type Repositories struct {
	Users  Users
	Admins Admins
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
