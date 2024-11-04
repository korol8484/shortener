package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestNewLogger(t *testing.T) {
	l, err := NewLogger(true)
	if err != nil {
		t.Fatal(err)
	}

	if l.Level() != zap.DebugLevel {
		t.Fatal("log not debug mode")
	}
}
