// Package model represents domain model. Every domain model type should have it's own file.
// It shouldn't depends on any other package in the application.
// It should only has domain model type and limited domain logic, in this example, validation logic. Because all other
// package depends on this package, the import of this package should be as small as possible.

package domain

import (
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	Id           string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserName     string    `bson:"username,omitempty" json:"username,omitempty"`
	FirstName    string    `bson:"firstname" json:"firstname"`
	LastName     string    `bson:"lastname" json:"lastname"`
	Email        string    `bson:"email,omitempty" json:"email,omitempty" `
	Mobile       string    `bson:"mobile" json:"mobile"`
	ShowReel     string    `bson:"showreel" json:"showreel"`
	ProfileImage string    `bson:"profileimage" json:"profileimage"`
	Bio          string    `bson:"bio" json:"bio"`
	Address      Address   `bson:"address" json:"address"`
	Pro          bool      `bson:"pro,omitempty" json:"pro,omitempty"`
	Created      time.Time `bson:"created,omitempty" json:"created,omitempty"`
}

// Render for All Responses
func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.NotNil.Error("invalid input")),
	)
}

// ValidatePersisted validate a user that has been persisted to database, basically Id is not empty
func (a User) ValidatePersisted() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required, validation.Length(3, 126)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Created, validation.Required),
	)
}
