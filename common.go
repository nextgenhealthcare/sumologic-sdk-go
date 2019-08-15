package sumologic

import (
	"errors"
)

// Error is returned by the API
type Error struct {
	Status  int    `json:"status"`
	ID      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrSourceNotFound is returned when a source doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if exists).
var ErrSourceNotFound = errors.New("Source not found")

// ErrAwsAuthenticationError is returned for authentication errors with AWS.
// Due to IAM's eventual consistency, it may be useful to retry.
var ErrAwsAuthenticationError = errors.New("Authentication Error with Sumo Logic")

type Filter struct {
	FilterType string `json:"filterType,omitempty"`
	Name       string `json:"name,omitempty"`
	Regexp     string `json:"regexp,omitempty"`
}
