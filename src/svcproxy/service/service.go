package service

var _ Service = &service{}

type service struct {
	proxies map[string]*Proxy
}

func NewService() *service {
	return &service{}
}

func (s *service) AddProxy(p *Proxy) error {
	return nil
}
