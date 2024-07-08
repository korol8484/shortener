package config

type App struct {
	Listen       string
	BaseShortURL string
}

func (a *App) GetBaseShortURL() string {
	return a.BaseShortURL
}
