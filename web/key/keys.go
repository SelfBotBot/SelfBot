package key

import "strings"

// Session keys.
const (
	SessionCore = "SelfBot"
)

// Context keys.
const (
	ContextUser     = "sb_user"
	ContextRedirect = "sb_redirectTo"
)

func IsContext(key string) bool {
	return strings.HasPrefix(key, "sb_")
}
