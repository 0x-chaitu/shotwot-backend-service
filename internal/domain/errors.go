package domain

import "errors"

var (
	ErrAccountNotFound         = errors.New("user doesn't exists")
	ErrNotFound                = errors.New("id doesn't exists")
	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
	ErrAccountAlreadyExists    = errors.New("user email already exists")
	ErrEmailPasswordInvalid    = errors.New("email or password is invalid")
	ErrNotAuthorized           = errors.New("user not authorized")
	ErrInvalidInput            = errors.New("invalid input")
)
