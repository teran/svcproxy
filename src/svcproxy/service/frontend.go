package service

// NewFrontend creates new Frontend instance
func NewFrontend(fqdn, httpHandler string, headers map[string]string) (*Frontend, error) {
	return &Frontend{
		FQDN:                fqdn,
		HTTPHandler:         httpHandler,
		ResponseHTTPHeaders: headers,
	}, nil
}
