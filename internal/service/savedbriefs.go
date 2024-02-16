package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
)

type SavedBriefsService struct {
	repo repository.SavedBriefs
}

func NewSavedBriefsService(repo repository.SavedBriefs) *SavedBriefsService {
	return &SavedBriefsService{
		repo: repo,
	}
}

func (b *SavedBriefsService) CreateOrUpdate(ctx context.Context, input *domain.SavedBriefInput) (*domain.SavedBriefRes, error) {

	savedBrief, err := b.repo.CreateOrUpdate(ctx, input.SavedBrief)
	if err != nil {
		return nil, err
	}
	return &domain.SavedBriefRes{
		SavedBrief: savedBrief,
	}, nil

}

func (b *SavedBriefsService) GetSavedBriefs(ctx context.Context, userId string) ([]*domain.SavedBrief, error) {
	return b.repo.GetSavedBriefs(ctx, userId)
}
