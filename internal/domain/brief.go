package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Brief struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string             `bson:"title,omitempty" json:"title,omitempty"`
	Images       []string           `bson:"images" json:"image"`
	Reward       int64              `bson:"reward,omitempty" json:"reward,omitempty"`
	Description  string             `bson:"description" json:"description"`
	Tags         []string           `bson:"tags" json:"tags"`
	ShotwotIdeas []*ShotwotIdeas    `bson:"shotwotIdeas" json:"shotwotIdeas"`
	References   []string           `bson:"references" json:"references"`
	Priority     int                `bson:"priority" json:"priority"`

	// timelapse, adventure, abstract
	Category    []string `bson:"category" json:"category"`
	Camera      []string `bson:"camera" json:"camera"`
	TechDetails []string `bson:"techDetails" json:"techDetails"`
	// 0 is indoor, 1 outdoor, 2 studio
	LightSetup []int `bson:"lightSetup" json:"lightSetup"`
	// audio, video, photo
	Type []string `bson:"type" json:"type"`

	CreatedBy string `bson:"createdby,omitempty" json:"createdBy,omitempty"`

	Created time.Time `bson:"created,omitempty" json:"created,omitempty"`

	IsActive  bool `bson:"isActive" json:"isActive"`
	Archieved bool `bson:"archieved" json:"archieved"`

	ActivatedOn *time.Time `bson:"activated_on,omitempty" json:"activatedOn,omitempty"`
	DeactivedOn *time.Time `bson:"deactivated_on,omitempty" json:"deactivatedOn,omitempty"`
}

type ShotwotIdeas struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}
