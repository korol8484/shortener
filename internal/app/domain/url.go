package domain

type URL struct {
	URL     string
	Alias   string
	Deleted bool
}

type BatchURL []*URL
