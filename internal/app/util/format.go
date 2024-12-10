package util

import "fmt"

// AddUrlToAlias - set url prefix for alias
func AddUrlToAlias(URL string) func(alias string) string {
	return func(alias string) string {
		return fmt.Sprintf("%s/%s", URL, alias)
	}
}
