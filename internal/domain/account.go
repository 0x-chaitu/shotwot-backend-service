// Package model represents domain model. Every domain model type should have it's own file.
// It shouldn't depends on any other package in the application.
// It should only has domain model type and limited domain logic, in this example, validation logic. Because all other
// package depends on this package, the import of this package should be as small as possible.

package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Account struct {
	Id           string    `bson:"_id" json:"id"`
	FirstName    string    `bson:"firstname" json:"firstname"`
	LastName     string    `bson:"lastname" json:"lastname"`
	Email        string    `bson:"email" json:"email" `
	Mobile       string    `bson:"mobile" json:"mobile"`
	ShowReel     string    `bson:"showreel" json:"showreel"`
	ProfileImage string    `bson:"profileimage" json:"profileimage"`
	Bio          string    `bson:"bio" json:"bio"`
	Address      Address   `bson:"address" json:"address"`
	Pro          bool      `bson:"pro" json:"pro"`
	Created      time.Time `bson:"created" json:"created"`
}

// ValidatePersisted validate a user that has been persisted to database, basically Id is not empty
func (a Account) ValidatePersisted() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required, validation.Length(3, 126)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Created, validation.Required),
	)
}
