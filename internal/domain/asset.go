package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Images = iota + 1
	Videos
	Audios
)

type Asset struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	AssetFile string             `bson:"assetFile,omitempty" json:"assetFile,omitempty"`
	UserId    string             `bson:"userId,omitempty" json:"userId,omitempty"`
	BriefId   primitive.ObjectID `bson:"briefId,omitempty" json:"briefId,omitempty"`

	AssetTitle   string   `bson:"assetTitle,omitempty" json:"assetTitle,omitempty"`
	Usage        string   `bson:"usage" json:"usage"`
	Category     []string `bson:"category" json:"category"`
	Description  string   `bson:"description" json:"description"`
	Tags         []string `bson:"tags" json:"tags"`
	RegionalTags []string `bson:"regionalTags" json:"regionalTags"`
	TimeOfDay    string   `bson:"timeOfDay" json:"timeOfDay"`
	Setting      string   `bson:"setting" json:"setting"`
	State        string   `bson:"state" json:"state"`
	Country      string   `bson:"country" json:"country"`
	Published    bool     `bson:"published" json:"published"`
	Archieved    bool     `bool:"archieved" json:"archieved"`
	Rating       *Rating  `bson:"rating" json:"rating"`
	Type         int      `bson:"type" json:"type"`

	Created time.Time `bson:"created,omitempty" json:"created,omitempty"`
	Updated time.Time `bson:"updated, omitempty" json:"updated"`
}

type AssetInput struct {
	*Asset
	File File `bson:"file" json:"file"`
}

type AssetRes struct {
	*Asset
	Url string `bson:"url" json:"url"`
}

type Rating struct {
	Current int `bson:"current" json:"current"`
	Total   int `bson:"total" json:"total"`
}
