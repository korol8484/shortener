package domain

// URL Base struct for app domain
type URL struct {
	URL     string
	Alias   string
	Deleted bool
}

// BatchURL Base collection struct for app domain
type BatchURL []*URL

// FormatAlias - format alias
func (u *URL) FormatAlias(f func(alias string) string) string {
	return f(u.Alias)
}
