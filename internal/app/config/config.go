package config

type App struct {
	Listen          string `env:"SERVER_ADDRESS"`
	BaseShortURL    string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func (a *App) GetBaseShortURL() string {
	return a.BaseShortURL
}

func (a *App) GetStoragePath() string {
	return a.FileStoragePath
}
