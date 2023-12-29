package domain

import "errors"

var (
	ErrAccountNotFound         = errors.New("user doesn't exists")
	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
	ErrAccountAlreadyExists    = errors.New("user with such email already exists")
)
