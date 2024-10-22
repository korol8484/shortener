package handlers

import (
	"net/http"
)

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
