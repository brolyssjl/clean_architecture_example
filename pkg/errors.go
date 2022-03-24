package pkg

import "errors"

var (
	ErrNotFound     = errors.New("error: Document not found")
	ErrNoContent    = errors.New("error: Document not found")
	ErrInvalidSlug  = errors.New("error: Invalid slug")
	ErrExists       = errors.New("error: Document already exists")
	ErrDatabase     = errors.New("error: Database error")
	ErrUnauthorized = errors.New("error: You are not allowed to perform this action")
	ErrForbidden    = errors.New("error: Access to this resource is forbidden")
)
