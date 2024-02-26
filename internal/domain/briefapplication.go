package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// status for application
const (
	Applied = iota + 1
	Shortlisted
	Rejected
	Accepted
)

type BriefApplication struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId  string             `bson:"userId,omitempty" json:"userId,omitempty"`
	Media   []string           `bson:"media,omitempty" json:"media,omitempty"`
	Reel    []string           `bson:"reel,omitempty" json:"reel,omitempty"`
	Note    string             `bson:"note,omitempty" json:"note,omitempty"`
	BriefId primitive.ObjectID `bson:"briefId,omitempty" json:"briefId,omitempty"`
	Created time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Time    time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	Opened  bool               `bson:"opened" json:"opened"`
	Status  int                `bson:"status,omitempty" json:"status,omitempty"`
	User    *struct {
		UserName string `bson:"username,omitempty" json:"username,omitempty"`
	} `bson:"user,omitempty" json:"user,omitempty"`
}

type UserBriefAppliedDetails struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId  string             `bson:"userId,omitempty" json:"userId,omitempty"`
	Media   []string           `bson:"media,omitempty" json:"media,omitempty"`
	Reel    []string           `bson:"reel,omitempty" json:"reel,omitempty"`
	Note    string             `bson:"note,omitempty" json:"note,omitempty"`
	BriefId primitive.ObjectID `bson:"briefId,omitempty" json:"briefId,omitempty"`
	Created time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Time    time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	Opened  bool               `bson:"opened" json:"opened"`
	Status  int                `bson:"status,omitempty" json:"status,omitempty"`
	User    *User              `bson:"user,omitempty" json:"user,omitempty"`
	Brief   *Brief             `bson:"brief,omitempty" json:"brief,omitempty"`
}
