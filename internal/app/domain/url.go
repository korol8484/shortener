package domain

type URL struct {
	URL   string
	Alias string
}

type BatchURL []*URL
