package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/helper"
)

type BriefsService struct {
	repo repository.Briefs
}

func NewBriefsService(repo repository.Briefs) *BriefsService {
	return &BriefsService{
		repo: repo,
	}
}

func (b *BriefsService) Create(ctx context.Context, input *domain.Brief) error {
	// validate input

	return b.repo.Create(ctx, input)

}

func (b *BriefsService) GetBriefs(ctx context.Context, predicate *helper.Predicate) ([]*domain.Brief, error) {
	return b.repo.GetBriefs(ctx, predicate)
}
