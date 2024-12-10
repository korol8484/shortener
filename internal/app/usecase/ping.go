package usecase

// PingDummy - Dummy interface
type PingDummy struct{}

// NewPingDummy - Dummy factory
func NewPingDummy() *PingDummy {
	return &PingDummy{}
}

// Ping - Dummy
func (p *PingDummy) Ping() error {
	return nil
}
