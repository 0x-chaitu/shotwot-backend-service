package domain

import (
	"net/http"
	"time"
)

type Admin struct {
	Id        string    `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string    `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName  string    `bson:"lastname" json:"lastname"`
	Email     string    `bson:"email,omitempty" json:"email,omitempty" `
	Mobile    string    `bson:"mobile" json:"mobile"`
	Role      int       `bson:"role,omitempty" json:"role,omitempty"`
	Created   time.Time `bson:"created,omitempty" json:"created,omitempty"`
}

// Render for All Responses
func (*Admin) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
