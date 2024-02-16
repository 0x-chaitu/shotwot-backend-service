package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavedBrief struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BriefId primitive.ObjectID `bson:"briefId,omitempty" json:"briefId,omitempty"`
	UserId  string             `bson:"userId,omitempty" json:"userId,omitempty"`
	Created time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Updated time.Time          `bson:"updated,omitempty" json:"updated,omitempty"`
	Status  bool               `bson:"status" json:"status"`
	User    User               `bson:"user,omitempty" json:"-"`
	Brief   Brief              `bson:"brief,omitempty" json:"brief"`
}

func NewSavedBrief() *SavedBrief {
	return &SavedBrief{
		Status: true,
	}
}

type SavedBriefInput struct {
	*SavedBrief
}

type SavedBriefRes struct {
	*SavedBrief
}
