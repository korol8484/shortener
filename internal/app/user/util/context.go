package util

import (
	"context"
)

type ctxKey string

const keyUserID ctxKey = "user_id"

func SetUserIDToCtx(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

func ReadUserIDFromCtx(ctx context.Context) int64 {
	userID, ok := ctx.Value(keyUserID).(int64)
	if !ok {
		return 0
	}

	return userID
}
