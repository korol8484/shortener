package util

import "fmt"

// AddURLToAlias - set url prefix for alias
func AddURLToAlias(URL string) func(alias string) string {
	return func(alias string) string {
		return fmt.Sprintf("%s/%s", URL, alias)
	}
}
