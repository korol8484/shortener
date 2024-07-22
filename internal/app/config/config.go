package config

type App struct {
	Listen       string `env:"SERVER_ADDRESS"`
	BaseShortURL string `env:"BASE_URL"`
}

func (a *App) GetBaseShortURL() string {
	return a.BaseShortURL
}
