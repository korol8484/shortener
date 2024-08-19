package util

import (
	"context"
)

type ctxKey string

const keyUserID ctxKey = "user_id"

func SetUserIdToCtx(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, keyUserID, userId)
}

func ReadUserIdFromCtx(ctx context.Context) int64 {
	userID, ok := ctx.Value(keyUserID).(int64)
	if !ok {
		return 0
	}

	return userID
}
