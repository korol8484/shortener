package storage

import (
	"context"
	"testing"
)

func TestMemoryStore_NewUser(t *testing.T) {
	store := NewMemoryStore()
	user, err := store.NewUser(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 {
		t.Fatal("user ID not valid, expexted 1")
	}
}
