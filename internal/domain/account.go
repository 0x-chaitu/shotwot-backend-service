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
	Id       string    `json:"uid"`
	Name     string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
}

// Validate validates a newly created user, which has not persisted to database yet, so Id is empty
func (a Account) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 126)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required, validation.Length(6, 72)),
	)
}

// ValidatePersisted validate a user that has been persisted to database, basically Id is not empty
func (a Account) ValidatePersisted() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 126)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required, validation.Length(6, 72)),
		validation.Field(&a.Created, validation.Required),
	)
}
