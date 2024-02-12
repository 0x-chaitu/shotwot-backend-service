package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"
	"shotwot_backend/pkg/s3"
	"time"
)

type BriefsService struct {
	repo repository.Briefs
	s3   *s3.S3Client
}

func NewBriefsService(repo repository.Briefs, s3 *s3.S3Client) *BriefsService {
	return &BriefsService{
		repo: repo,
		s3:   s3,
	}
}

func (b *BriefsService) Create(ctx context.Context, input *domain.BriefInput) (*domain.BriefRes, error) {
	var images = []string{}
	var urls = []string{}
	timeKey := time.Now().Format("20060102150405")
	for _, file := range input.Files {
		key := "app/brief/" + timeKey + "/" + file.Name
		url, err := b.s3.PresignedUrl(key, file.Filetype)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url.URL)
		images = append(images, key)
	}
	var cardUrl string
	if input.CardFile != nil {
		key := "app/brief/" + timeKey + "/" + input.CardFile.Name
		url, err := b.s3.PresignedUrl(key, input.CardFile.Filetype)
		if err != nil {
			return nil, err
		}
		input.CardImage = key
		cardUrl = url.URL
	}

	input.Brief.Images = images
	brief, err := b.repo.Create(ctx, input.Brief)
	if err != nil {
		return nil, err
	}
	return &domain.BriefRes{
		Brief:   brief,
		Urls:    urls,
		CardUrl: cardUrl,
	}, nil

}

func (b *BriefsService) GetBriefs(ctx context.Context, predicate *helper.BriefPredicate) ([]*domain.Brief, error) {
	return b.repo.GetBriefs(ctx, predicate)
}

func (b *BriefsService) DeleteBrief(ctx context.Context, id string) error {
	return b.repo.DeleteBrief(ctx, id)
}

func (b *BriefsService) Update(ctx context.Context, input *domain.BriefInput) (*domain.Brief, error) {
	brief, err := b.repo.Update(ctx, input.Brief)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return brief, nil
}
