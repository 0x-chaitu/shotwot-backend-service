package domain

import "time"

type Video struct {
	Id    string `bson:"_id,omitempty" json:"id,omitempty"`
	Link  string `bson:"link" json:"link"`
	Title string `bson:"title" json:"title"`
}

type Playlist struct {
	Id             string    `bson:"_id,omitempty" json:"id,omitempty"`
	Title          string    `bson:"title" json:"title"`
	About          string    `bson:"about" json:"about"`
	ThumbnailImage string    `bson:"thumbnailImage" json:"thumbnailImage"`
	Videos         []Video   `bson:"videos,omitempty" json:"videos,omitempty"`
	IsActive       bool      `bson:"isActive" json:"isActive"`
	CreatedBy      string    `bson:"createdby,omitempty" json:"createdBy,omitempty"`
	Created        time.Time `bson:"created,omitempty" json:"created,omitempty"`
}

type PlaylistInput struct {
	*Playlist
}

type PlaylistResp struct {
	*Playlist
}
