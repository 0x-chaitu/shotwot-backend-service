package domain

import "time"

type Admin struct {
	Id        string    `bson:"_id" json:"id"`
	FirstName string    `bson:"firstname" json:"firstname"`
	LastName  string    `bson:"lastname" json:"lastname"`
	Email     string    `bson:"email" json:"email" `
	Mobile    string    `bson:"mobile" json:"mobile"`
	Role      string    `bson:"role" json:"role"`
	Created   time.Time `bson:"created" json:"created"`
}
