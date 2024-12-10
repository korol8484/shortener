package usecase

import "testing"

func TestNewPingDummy(t *testing.T) {
	dp := NewPingDummy()
	err := dp.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
