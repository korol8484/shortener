package util

import (
	"context"
)

type ctxKey string

// keyUserID - key to set\read user from context
const keyUserID ctxKey = "user_id"

// SetUserIDToCtx - add user ID in context
func SetUserIDToCtx(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

// ReadUserIDFromCtx - read user ID from context
func ReadUserIDFromCtx(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(keyUserID).(int64)
	if !ok {
		return 0, false
	}

	return userID, true
}
