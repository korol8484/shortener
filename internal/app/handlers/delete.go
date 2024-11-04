package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/korol8484/shortener/internal/app/user/util"
)

type batchItem struct {
	aliases []string
	user    int64
}

// Delete BatchDelete Handler
type Delete struct {
	store     Store
	batchChan chan batchItem
	closeChan chan struct{}
	logger    *zap.Logger
	batchSize int
}

// NewDelete Factory for BatchDelete Handler
func NewDelete(store Store, logger *zap.Logger) (*Delete, error) {
	d := &Delete{
		store:     store,
		batchChan: make(chan batchItem, 1024),
		closeChan: make(chan struct{}),
		logger:    logger,
		batchSize: 500,
	}

	for i := 0; i < 2; i++ {
		go d.process()
	}

	return d, nil
}

// BatchDelete Handler for a collection of delete user shorten URLs
// Accepts input json:
// ["cbi7jn", "dyifOs"]
// Returns: Http status Accepted
func (d *Delete) BatchDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := util.ReadUserIDFromCtx(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var aliases []string

	if err = json.Unmarshal(body, &aliases); err != nil {
		http.Error(w, "can't parse json", http.StatusBadRequest)
		return
	}

	go d.add(aliases, userID)

	w.WriteHeader(http.StatusAccepted)
}

// Close - close resources
func (d *Delete) Close() {
	close(d.closeChan)
	close(d.batchChan)
}

func (d *Delete) add(aliases []string, userID int64) {
	for i := 0; i < len(aliases); i += d.batchSize {
		end := i + d.batchSize
		if end > len(aliases) {
			end = len(aliases)
		}

		d.batchChan <- batchItem{
			aliases: aliases[i:end],
			user:    userID,
		}
	}
}

func (d *Delete) process() {
	for {
		select {
		case batch, ok := <-d.batchChan:
			if !ok {
				return
			}

			if err := d.store.BatchDelete(context.Background(), batch.aliases, batch.user); err != nil {
				d.logger.Error(
					"can't delete bach",
					zap.Int64("userId", batch.user),
					zap.Error(err),
				)
			}
		case <-d.closeChan:
			d.logger.Info("close delete worker")
			return
		}
	}
}
