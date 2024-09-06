package config

type App struct {
	Listen          string `env:"SERVER_ADDRESS"`
	BaseShortURL    string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DBDsn           string `env:"DATABASE_DSN"`
}

func (a *App) GetBaseShortURL() string {
	return a.BaseShortURL
}

func (a *App) GetStoragePath() string {
	return a.FileStoragePath
}

// GetDsn - postgresql://postgres:postgres@localhost:5432/short
func (a *App) GetDsn() string {
	return a.DBDsn
}
