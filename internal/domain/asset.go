package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Asset struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	AssetFile []string           `bson:"assetFile,omitempty" json:"assetFile,omitempty"`
	UserId    string             `bson:"userId,omitempty" json:"userId,omitempty"`
	BriefId   primitive.ObjectID `bson:"briefId,omitempty" json:"briefId,omitempty"`

	Created time.Time `bson:"created,omitempty" json:"created,omitempty"`
}

type AssetInput struct {
	*Asset
	Files []File `bson:"files" json:"files"`
}

type AssetRes struct {
	*Asset
	Urls []string `bson:"urls" json:"urls"`
}
