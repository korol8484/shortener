package util

import (
	"context"
	"testing"
)

func TestReadUserIDFromCtx(t *testing.T) {
	cxt := context.Background()

	ctxWithUser := SetUserIDToCtx(cxt, 1)
	userID, ok := ReadUserIDFromCtx(ctxWithUser)
	if !ok {
		t.Fatal("user ID not set in context")
	}

	if userID != 1 {
		t.Fatal("user ID not equal")
	}

	_, notOk := ReadUserIDFromCtx(cxt)
	if notOk != false {
		t.Fatal("user can't be found in context")
	}
}
