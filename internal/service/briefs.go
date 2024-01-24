package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"
)

type BriefsService struct {
	repo repository.Briefs
}

func NewBriefsService(repo repository.Briefs) *BriefsService {
	return &BriefsService{
		repo: repo,
	}
}

func (b *BriefsService) Create(ctx context.Context, input *domain.Brief) (*domain.Brief, error) {
	// validate input

	return b.repo.Create(ctx, input)

}

func (b *BriefsService) GetBriefs(ctx context.Context, predicate *helper.BriefPredicate) ([]*domain.Brief, error) {
	return b.repo.GetBriefs(ctx, predicate)
}

func (b *BriefsService) DeleteBrief(ctx context.Context, id string) error {
	return b.repo.DeleteBrief(ctx, id)
}

func (b *BriefsService) Update(ctx context.Context, input *domain.Brief) (*domain.Brief, error) {
	brief, err := b.repo.Update(ctx, input)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return brief, nil
}
