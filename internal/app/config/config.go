package config

// App application configuration
type App struct {
	// Listen host:port on which web service will operate
	Listen string `env:"SERVER_ADDRESS"`
	// BaseShortURL HTTP domain append to short URL
	BaseShortURL string `env:"BASE_URL"`
	// FileStoragePath Path to file database
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// DBDsn Database connection string
	DBDsn string `env:"DATABASE_DSN"`
}

// GetBaseShortURL return HTTP domain, append to short URL
func (a *App) GetBaseShortURL() string {
	return a.BaseShortURL
}

// GetStoragePath Path to file database
func (a *App) GetStoragePath() string {
	return a.FileStoragePath
}

// GetDsn Database connection string:
// Example: postgresql://postgres:postgres@localhost:5432/short
func (a *App) GetDsn() string {
	return a.DBDsn
}
