package storerrros

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("user alredy exists")
	ErrUserNoExist     = errors.New("user does not exists")

	ErrBookNoExist    = errors.New("book does not exists")
	ErrEmptyBooksList = errors.New("empty books list")
)
