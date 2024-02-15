package repository

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/helper"

	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user *domain.User) error

	Update(ctx context.Context, user *domain.User) (*domain.User, error)

	Get(ctx context.Context, id string) (*domain.User, error)

	GetOrCreate(ctx context.Context, user *domain.User) (*domain.User, error)

	GetUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	TotalUsers(ctx context.Context) (int64, error)

	Download(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)

	SearchUsers(ctx context.Context, predicate *helper.UsersPredicate) ([]*domain.User, error)
}

type Admins interface {
	Create(ctx context.Context, admin *domain.Admin) error

	Update(ctx context.Context, admin *domain.Admin) (*domain.Admin, error)

	Get(ctx context.Context, id string) (*domain.Admin, error)

	GetAdmins(ctx context.Context) ([]*domain.Admin, error)
}

type Briefs interface {
	Create(ctx context.Context, brief *domain.Brief) (*domain.Brief, error)

	Update(ctx context.Context, brief *domain.Brief) (*domain.Brief, error)

	GetBrief(ctx context.Context, id string) (*domain.Brief, error)

	GetBriefs(ctx context.Context, predicate *helper.BriefPredicate) ([]*domain.Brief, error)

	DeleteBrief(ctx context.Context, id string) error
}
type BriefApplications interface {
	Create(ctx context.Context, briefapplication *domain.BriefApplication) (*domain.BriefApplication, error)

	GetBriefApplications(ctx context.Context, id string) ([]*domain.BriefApplication, error)
}

type Repositories struct {
	Users             Users
	Admins            Admins
	Briefs            Briefs
	BriefApplications BriefApplications
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users:             NewUsersRepo(db),
		Admins:            NewAdminsRepo(db),
		Briefs:            NewBriefsRepo(db),
		BriefApplications: NewBriefApplicationsRepo(db),
	}
}
