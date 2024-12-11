package usecase

import (
	"context"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/util"
	"go.uber.org/zap"
	"hash/fnv"
	"math/rand"
	"net/url"
)

// Config Return HTTP domain append to short URL
type Config interface {
	GetBaseShortURL() string
}

// Pingable check service status contract
type Pingable interface {
	Ping() error
}

// Store Repository Interface
type Store interface {
	Add(ctx context.Context, ent *domain.URL, user *domain.User) error
	Read(ctx context.Context, alias string) (*domain.URL, error)
	AddBatch(ctx context.Context, batch domain.BatchURL, user *domain.User) error
	ReadUserURL(ctx context.Context, user *domain.User) (domain.BatchURL, error)
	BatchDelete(ctx context.Context, aliases []string, userID int64) error
	LoadStats(ctx context.Context) (*domain.StatsModel, error)
	Close() error
}

type batchItem struct {
	aliases []string
	user    int64
}

// Usecase - base service usage
type Usecase struct {
	health Pingable
	store  Store
	cfg    Config
	logger *zap.Logger

	batchChan chan batchItem
	closeChan chan struct{}
	batchSize int
}

// NewUsecase factory
func NewUsecase(cfg Config, store Store, health Pingable, logger *zap.Logger) *Usecase {
	u := &Usecase{
		health:    health,
		store:     store,
		cfg:       cfg,
		logger:    logger,
		batchChan: make(chan batchItem, 1024),
		closeChan: make(chan struct{}),
		batchSize: 100,
	}

	for i := 0; i < 2; i++ {
		go u.process()
	}

	return u
}

// Ping check service HTTP handler
func (u *Usecase) Ping() bool {
	if err := u.health.Ping(); err != nil {
		return false
	}

	return true
}

// CreateURL create shorten url for user
func (u *Usecase) CreateURL(ctx context.Context, URL string, userID int64) (*domain.URL, error) {
	ent, err := u.GenerateURL(URL)
	if err != nil {
		return nil, err
	}

	if err = u.store.Add(ctx, ent, &domain.User{ID: userID}); err != nil {
		return ent, err
	}

	return ent, nil
}

// LoadByAlias load from store by shorten url
func (u *Usecase) LoadByAlias(ctx context.Context, alias string) (*domain.URL, error) {
	return u.store.Read(ctx, alias)
}

// LoadAllUserURL load all url by user
func (u *Usecase) LoadAllUserURL(ctx context.Context, userID int64) (domain.BatchURL, error) {
	return u.store.ReadUserURL(ctx, &domain.User{ID: userID})
}

// GetStats load statistic
func (u *Usecase) GetStats(ctx context.Context) (*domain.StatsModel, error) {
	return u.store.LoadStats(ctx)
}

// AddBatch add batch url to store
func (u *Usecase) AddBatch(ctx context.Context, batch domain.BatchURL, userID int64) error {
	return u.store.AddBatch(ctx, batch, &domain.User{ID: userID})
}

// FormatAlias add to alias base short url
func (u *Usecase) FormatAlias(URL *domain.URL) string {
	return URL.FormatAlias(util.AddURLToAlias(u.cfg.GetBaseShortURL()))
}

// GenerateURL generate short alias by url
func (u *Usecase) GenerateURL(URL string) (*domain.URL, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	ent := &domain.URL{
		URL:   parsedURL.String(),
		Alias: GenAlias(6, parsedURL.String()),
	}

	return ent, nil
}

// AddToDelete add task to delete url for user
func (u *Usecase) AddToDelete(aliases []string, userID int64) {
	go func() {
		for i := 0; i < len(aliases); i += u.batchSize {
			end := i + u.batchSize
			if end > len(aliases) {
				end = len(aliases)
			}

			u.batchChan <- batchItem{
				aliases: aliases[i:end],
				user:    userID,
			}
		}
	}()
}

func (u *Usecase) process() {
	for {
		select {
		case batch, ok := <-u.batchChan:
			if !ok {
				return
			}

			if err := u.store.BatchDelete(context.Background(), batch.aliases, batch.user); err != nil {
				u.logger.Error(
					"can't delete bach",
					zap.Int64("userId", batch.user),
					zap.Error(err),
				)
			}
		case <-u.closeChan:
			u.logger.Info("close delete worker")
			return
		}
	}
}

// Close - close resources
func (u *Usecase) Close() {
	close(u.closeChan)
	close(u.batchChan)

	_ = u.store.Close()
}

// GenAlias - Create alias length n as a hash of the string
func GenAlias(keyLen int, shortURL string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	h := fnv.New64()
	h.Write([]byte(shortURL))

	r := rand.New(rand.NewSource(int64(h.Sum64())))

	keyMap := make([]byte, keyLen)
	for i := range keyMap {
		keyMap[i] = charset[r.Intn(len(charset))]
	}

	return string(keyMap)
}
