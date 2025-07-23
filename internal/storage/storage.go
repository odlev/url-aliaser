// Package storage is a nice package
package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExist = errors.New("url already exist")
	ErrAliasNotFound = errors.New("alias not found")
)
