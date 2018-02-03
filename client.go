// Package sumologic is a wrapper for the Sumo Logic API.
package sumologic

import (
	"errors"
	"net/url"
)

// Client communicates with the Sumo Logic API.
type Client struct {
	AuthToken   string
	EndpointURL *url.URL
}

// ErrClientAuthenticationError is returned for authentication errors with the API.
var ErrClientAuthenticationError = errors.New("Authentication Error with Sumo Logic")

// NewClient returns a new sumologic.Client for accessing the Sumo Logic API.
func NewClient(authToken, defaultEndpointURL string) (*Client, error) {
	s := &Client{
		AuthToken: authToken,
	}
	endpointURL, err := url.Parse(defaultEndpointURL)
	if err != nil {
		return nil, err
	}
	s.EndpointURL = endpointURL
	return s, nil
}
