package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Brief struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string             `bson:"title" json:"title"`
	CardImage    string             `bson:"cardImage" json:"cardImage"`
	Images       []string           `bson:"images" json:"images"`
	Reward       int64              `bson:"reward" json:"reward"`
	Description  string             `bson:"description" json:"description"`
	Tags         []string           `bson:"tags" json:"tags"`
	ShotwotIdeas []*ShotwotIdeas    `bson:"shotwotIdeas" json:"shotwotIdeas"`
	Priority     int                `bson:"priority" json:"priority"`

	// timelapse, adventure, abstract
	Camera     []string `bson:"camera" json:"camera"`
	Resolution string   `bson:"resolution" json:"resolution"`
	FrameRate  string   `bson:"frameRate" json:"frameRate"`
	Duration   string   `bson:"duration" json:"duration"`
	// 0 is indoor, 1 outdoor, 2 studio
	LightSetup []string `bson:"lightSetup" json:"lightSetup"`
	// audio, video, photo
	Type string `bson:"type" json:"type"`

	CreatedBy string `bson:"createdby,omitempty" json:"createdBy,omitempty"`

	Created time.Time `bson:"created,omitempty" json:"created,omitempty"`
	Expiry  time.Time `bson:"expiry,omitempty" json:"expiry,omitempty"`

	IsActive    bool `bson:"isActive" json:"isActive"`
	TotalAssets int  `bson:"assets" json:"assets"`

	// ActivatedOn *time.Time `bson:"activated_on,omitempty" json:"activatedOn,omitempty"`
	// DeactivedOn *time.Time `bson:"deactivated_on,omitempty" json:"deactivatedOn,omitempty"`
}

type ShotwotIdeas struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}

type BriefInput struct {
	*Brief
	Files    []File `bson:"files" json:"files"`
	CardFile *File  `bson:"file" json:"file"`
}

type BriefRes struct {
	*Brief
	Urls    []string `bson:"urls" json:"urls"`
	CardUrl string   `bson:"url" json:"url"`
}
