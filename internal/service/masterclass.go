package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/s3"
	"time"
)

type MasterClassService struct {
	repo repository.MasterClass
	s3   *s3.S3Client
}

func NewMasterClassService(repo repository.MasterClass, s3 *s3.S3Client) *MasterClassService {
	return &MasterClassService{
		repo: repo,
		s3:   s3,
	}
}

func (b *MasterClassService) CreatePlaylist(ctx context.Context, input *domain.PlaylistInput) (*domain.PlaylistResp, error) {
	var videos []domain.Video
	timeKey := time.Now().Format("20060102150405")
	for _, video := range input.Videos {
		key := "app/masterclass/playlist/" + timeKey + "/" + video.Title // Assuming the video title is unique for the key
		if _, err := b.s3.PresignedUrl(key, video.Link); err != nil {    // Assuming 'Link' is the video URL
			return nil, err
		}
		video.Id = key
		videos = append(videos, video)
	}

	playlist := &domain.Playlist{
		Title:          input.Title,
		About:          input.About,
		ThumbnailImage: input.ThumbnailImage,
		Videos:         videos,
		IsActive:       input.IsActive,
	}

	createdPlaylist, err := b.repo.CreatePlaylist(ctx, playlist)
	if err != nil {
		return nil, err
	}

	return &domain.PlaylistResp{
		Playlist: createdPlaylist,
	}, nil
}

func (b *MasterClassService) GetPlaylists(ctx context.Context, predicate *helper.PlaylistPredicate) ([]*domain.Playlist, error) {
	return b.repo.GetPlaylists(ctx, predicate)
}
