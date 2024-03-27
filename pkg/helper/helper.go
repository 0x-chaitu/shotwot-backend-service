package helper

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	Ascending  = 1
	Descending = -1
)

func TODoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

type BriefPredicate struct {
	Order     int    `json:"order,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	IsActive  *bool  `json:"isActive,omitempty"`
	ByDate    string `json:"byDate,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`

	Type    []string `json:"type,omitempty" `
	Expiry  int      `json:"expiry,omitempty"`
	RewardL int      `json:"rewardL,omitempty"`
	RewardG int      `json:"rewardG,omitempty"`
	Skip    int64    `json:"skip,omitempty"`
}

type UsersPredicate struct {
	Order     int       `json:"order,omitempty"`
	StartDate time.Time `json:"created,omitempty"`
	EndDate   time.Time `json:"end,omitempty"`
	Skip      int       `json:"skip,omitempty"`
	Key       string    `json:"key,omitempty"`
}

type BriefApplicationsPredicate struct {
	Order     int    `json:"order,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	IsActive  *bool  `json:"isActive,omitempty"`
	ByDate    string `json:"byDate,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

type PlaylistPredicate struct {
	Order     int    `json:"order,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	IsActive  *bool  `json:"isActive,omitempty"`
	ByDate    string `json:"byDate,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}
