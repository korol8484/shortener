package handlers

import (
	"context"
	"encoding/json"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/user/util"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type batchItem struct {
	aliases []string
	user    *domain.User
}

type Delete struct {
	store     Store
	batchChan chan batchItem
	closeChan chan struct{}
	logger    *zap.Logger
	batchSize int
}

func NewDelete(store Store, logger *zap.Logger) (*Delete, error) {
	d := &Delete{
		store:     store,
		batchChan: make(chan batchItem, 10),
		closeChan: make(chan struct{}),
		logger:    logger,
		batchSize: 500,
	}

	for i := 0; i < 2; i++ {
		go d.process()
	}

	return d, nil
}

func (d *Delete) BatchDelete(w http.ResponseWriter, r *http.Request) {
	userID := util.ReadUserIDFromCtx(r.Context())

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

	go d.add(aliases, &domain.User{ID: userID})

	w.WriteHeader(http.StatusAccepted)
}

func (d *Delete) Close() {
	close(d.closeChan)
}

func (d *Delete) add(aliases []string, user *domain.User) {
	for i := 0; i < len(aliases); i += d.batchSize {
		end := i + d.batchSize
		if end > len(aliases) {
			end = len(aliases)
		}

		d.batchChan <- batchItem{
			aliases: aliases[i:end],
			user:    user,
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
					zap.Int64("userId", batch.user.ID),
					zap.Error(err),
				)
			}
		case <-d.closeChan:
			close(d.batchChan)
			d.logger.Info("close delete worker")
			return
		}
	}
}
