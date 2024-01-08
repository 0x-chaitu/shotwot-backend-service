package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Brief struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title      string             `bson:"title,omitempty" json:"title,omitempty"`
	Image      string             `bson:"image" json:"image"`
	Reward     int64              `bson:"reward,omitempty" json:"reward,omitempty"`
	Summary    string             `bson:"summary" json:"summary"`
	Resolution string             `bson:"resolution" json:"resolution"`

	IsActive  bool `bson:"is_active" json:"is_active"`
	Archieved bool `bson:"archieved" json:"archieved"`

	ActivatedOn time.Time `bson:"activated_on,omitempty" json:"activatedOn,omitempty"`
	DeactivedOn time.Time `bson:"deactivated_on,omitempty" json:"deactivatedOn,omitempty"`

	// timelapse, adventure, abstract
	Category string `bson:"category" json:"category"`

	// audio, video, photo
	Type string `bson:"type" json:"type"`

	CreatedBy string    `bson:"createdby,omitempty" json:"createdBy,omitempty"`
	Created   time.Time `bson:"created,omitempty" json:"created,omitempty"`
}
