package domain

type Video struct {
	Id    string `bson:"_id,omitempty" json:"id,omitempty"`
	Link  string `bson:"link" json:"link"`
	Title string `bson:"title" json:"title"`
}

type Playlist struct {
	Id             string  `bson:"_id,omitempty" json:"id,omitempty"`
	Title          string  `bson:"title" json:"title"`
	About          string  `bson:"about" json:"about"`
	ThumbnailImage string  `bson:"thumbnailImage" json:"thumbnailImage"`
	Videos         []Video `bson:"videos,omitempty" json:"videos,omitempty"`
}
