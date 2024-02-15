package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BriefApplication struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId  string             `bson:"userId,omitempty" json:"userId,omitempty"`
	BriefId primitive.ObjectID `bson:"briefId,omitempty" json:"briefId,omitempty"`
	Created time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Time    time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	Status  string             `bson:"status,omitempty" json:"status,omitempty"`
	User    User               `bson:"user,omitempty" json:"user,omitempty"`
	Brief   Brief              `bson:"brief,omitempty" json:"brief,omitempty"`
}

type BriefApplicationInput struct {
	*BriefApplication
}

type BriefApplicationRes struct {
	*BriefApplication
}
