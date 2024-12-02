package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net"
	"net/http"
)

// StatsCfg - config interface
type StatsCfg interface {
	GetTrustedSubnet() string
}

// Stats - service handler
type Stats struct {
	logger  *zap.Logger
	ipNet   *net.IPNet
	storage Store
}

// NewStats - factory
func NewStats(cfg StatsCfg, logger *zap.Logger, storage Store) (*Stats, error) {
	var ipNet *net.IPNet
	var err error

	if cfg.GetTrustedSubnet() != "" {
		_, ipNet, err = net.ParseCIDR(cfg.GetTrustedSubnet())
		if err != nil {
			return nil, err
		}
	}

	return &Stats{
		logger:  logger,
		ipNet:   ipNet,
		storage: storage,
	}, nil
}

func (s *Stats) handle(w http.ResponseWriter, r *http.Request) {
	if s.ipNet == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ipStr := r.Header.Get("X-Real-IP")

	if ipStr == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !s.ipNet.Contains(ip) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stat, err := s.storage.LoadStats(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(stat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", mimeJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
