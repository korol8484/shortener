package handlers

import (
	"net/http"
)

// Pingable check service status contract
type Pingable interface {
	Ping() error
}

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

// Ping check service HTTP handler
func Ping(p Pingable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
