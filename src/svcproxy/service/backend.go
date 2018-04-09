package service

import (
	"net/url"
)

// NewBackend creates new Backend instance
func NewBackend(address string) (*Backend, error) {
	backend, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL: backend,
	}, nil
}
