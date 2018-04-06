package service

// Service interface
type Service interface {
	AddProxy(*Proxy) error
}

type Proxy struct {
	frontend *Frontend
	backend  *Backend
}

// Frontend type
type Frontend struct{}

// Backend type
type Backend struct{}
